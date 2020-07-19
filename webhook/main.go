package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.POST("/", func(c *gin.Context) {
		c.JSON(202, gin.H{})

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(c.Request.Body)

		gitRepo, _ := jsonparser.GetString(buf.Bytes(), "repository", "full_name")
		appName, _ := jsonparser.GetString(buf.Bytes(), "repository", "name")
		gitRevision, _ := jsonparser.GetString(buf.Bytes(), "push", "changes", "[0]", "new", "name")

		if gitRevision == "" {
			return
		}
		fmt.Println("repo: ", gitRepo)
		fmt.Println("ref: ", gitRevision)
		fmt.Println("name: ", appName)

		gitRepoName := strings.Split(gitRepo, "/")[1]
		fullGitRepo := "git@bitbucket.org:" + gitRepo + ".git"
		// timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

		// argoFilename := "argo" + timestamp + ".yml"
		argoFilename := "argo_.yml"

		items := []string{"staging", "master"}
		_, found := Find(items, gitRevision)
		if !found {
			return
		}


		commandin := "argo submit " + argoFilename + " -p app-name=" + gitRepoName + " -p repo=" + fullGitRepo + " -p ref=" + gitRevision
		fmt.Println(commandin)

		commandOutput, err := exec.Command("sh", "-c", commandin).CombinedOutput()
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Accepted webhook request, did NOT start Argo workflow: git_repo=%q,git_revision=%q, because of: %q\n", gitRepo, gitRevision, string(err.Error()))
		} else {
			fmt.Printf("Accepted webhook request, started Argo workflow: git_repo=%q,git_revision=%q, with message: %q\n", gitRepo, gitRevision, string(commandOutput))
		}
	})
	_ = r.Run(":3000")
}


func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}