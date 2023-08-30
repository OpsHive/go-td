package helpers

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
)

func CommitAndPush(yamlBytes []byte, fileName string) {

	gitUser := os.Getenv("GIT_USERNAME")
	gitPass := os.Getenv("GIT_TOKEN")
	gitrepo := os.Getenv("GIT_HELM_REPO_URL")

	encodedGitUser := url.QueryEscape(gitUser)
	encodedGitPass := url.QueryEscape(gitPass)
	repoURL := fmt.Sprintf("https://%s:%s@%s", encodedGitUser, encodedGitPass, gitrepo)
	//repoURL := "https://boot:tripon-UBw_QsEvb78VaS8nkDFb@git.tripon.io/tripon/helm-charts.git"

	repoDir := "./tmp" // Set this to your desired directory

	if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
		fmt.Println("Failed to create repo directory")
		return
	}
	gitCmd := exec.Command("git", "clone", repoURL, repoDir)

	if err := gitCmd.Run(); err != nil {
		fmt.Println("Failed to clone repository")
		return
	}
	gitConfigCmd := exec.Command("git", "config", "--global", "credential.helper", "store")
	if err := gitConfigCmd.Run(); err != nil {
		// Handle the error
		return
	}
	// Write the YAML content to a temporary file
	filename := fmt.Sprintf("%v.yaml", fileName)
	tmpfilePath := "./tmp/" + filename
	tmpfile, err := os.Create(tmpfilePath)
	if err != nil {
		fmt.Println("Failed to create temporary file", err.Error())
		return
	}
	if _, err := tmpfile.Write(yamlBytes); err != nil {
		fmt.Println("Failed to write to temporary file")
		return
	}

	// Copy the temporary file to the repository directory

	// Stage and commit the changes
	gitAddCmd := exec.Command("git", "add", ".")
	gitAddCmd.Dir = repoDir
	if err := gitAddCmd.Run(); err != nil {
		fmt.Println("Failed to stage changes")
		return
	}

	gitCommitCmd := exec.Command("git", "commit", "-m", "Adding knative service YAML")
	gitCommitCmd.Dir = repoDir
	if err := gitCommitCmd.Run(); err != nil {
		fmt.Println("Failed to commit changes")
		return
	}

	// Set up Git credentials

	// Push the changes
	gitPushCmd := exec.Command("git", "push", "--set-upstream", "origin", "main")
	gitPushCmd.Dir = repoDir
	var output bytes.Buffer
	gitPushCmd.Stdout = &output
	gitPushCmd.Stderr = &output
	if err := gitPushCmd.Run(); err != nil {
		fmt.Println("Failed to push changes", output.String())
		return
	}

}
