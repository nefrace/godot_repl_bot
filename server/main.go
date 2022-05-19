package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CodeStruct struct {
	Code []string
	Id   string
}

func main() {
	r := gin.Default()
	var godot_path string
	godot_path, exists := os.LookupEnv("GODOT")
	if !exists {
		godot_path = "godot"
	}
	scripts_path, exists := os.LookupEnv("SCRIPTS")
	if !exists {
		scripts_path = "./"
	}
	r.POST("/run", func(c *gin.Context) {
		code := c.DefaultPostForm("code", "")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "Empty body"})
			return
		}
		lines := strings.Split(code, "\n")
		log.Printf("%v", lines)
		file_uuid := uuid.New()
		file_uuid_clean := strings.Replace(file_uuid.String(), "-", "", -1)
		filename := scripts_path + file_uuid_clean + ".gd"
		// err := os.WriteFile(filename, code_bytes, 0644)
		tmpl, err := template.ParseFiles("template")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not parse template", "error": err.Error()})
			return
		}
		file, err := os.Create(filename)
		defer file.Close()
		defer os.Remove(filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not create file", "error": err.Error()})
			return
		}
		err = tmpl.Execute(file, CodeStruct{Code: lines, Id: file_uuid_clean})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "could not write file", "error": err.Error()})
			return
		}
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, godot_path, "-s", filename, "--no-window")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		out_str := stdout.String()
		err_str := stderr.String()
		separator := "=== " + file_uuid_clean + " ==="
		out_sep := strings.Split(out_str, separator)
		err_sep := strings.Split(err_str, separator)
		if len(out_sep) == 1 {
			out_sep = append(out_sep, "")
		}
		if len(err_sep) == 1 {
			err_sep = append(err_sep, "")
		}
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{"status": "script timeout", "stdout": out_sep[1], "stderr": err_sep[1]})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}

		out_separated := strings.Split(out_str, separator)
		if len(out_separated) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "runtime error", "stdout": out_str, "stderr": err_str})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": out_separated[1]})
	})
	log.Fatal(r.Run(":8080").Error())
}
