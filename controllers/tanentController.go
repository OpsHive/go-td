package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opshive/go-td/controllers/helpers"
	"github.com/opshive/go-td/models"
	"github.com/opshive/go-td/pkg/clients"
	"gopkg.in/yaml.v2"
)

// @Summary Deploy Tenant
// @Description Deploy a new Tenant by providing a JSON object with the "name" property.
// @Tags Tenant
// @Accept json
// @Produce json
// @Param name body string true "Tenant Name"
// @Success 200 {string} string "Tenant applied successfully"
// @Router /tenantCreate [post]
func TenantDeploy(c *gin.Context) {
	var userInput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	// implement redis cache

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service := models.Service{
		APIVersion: "serving.knative.dev/v1",
		Kind:       "Service",
		Metadata: models.Metadata{
			Name:      userInput.Name,
			Namespace: "default",
			Labels: map[string]string{
				"app": "secret",
			},
		},
		Spec: models.Spec{
			Template: models.Template{
				Spec: models.SpecTemplate{
					ImagePullSecrets: []models.ImagePullSecret{
						{Name: "tripon"},
					},
					Containers: []models.Container{
						{
							Image: "ghcr.io/knative/helloworld-go:latest",
							Ports: []models.Port{
								{ContainerPort: 8080},
							},
							Env: []models.EnvVar{
								{
									Name:  "TARGET",
									Value: "World",
								},
							},
						},
					},
				},
			},
		},
	}

	yamlBytes, err := yaml.Marshal(service)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate YAML"})
		return
	}
	fileName := userInput.Name

	helpers.CommitAndPush(yamlBytes, fileName)

	// repoURL := "https://boot:tripon-UBw_QsEvb78VaS8nkDFb@git.tripon.io/tripon/helm-charts.git"
	// repoDir := "./tmp" // Set this to your desired directory
	// defer os.RemoveAll(repoDir)
	// if err := os.MkdirAll(repoDir, os.ModePerm); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repo directory"})
	// 	return
	// }
	// gitCmd := exec.Command("git", "clone", repoURL, repoDir)

	// if err := gitCmd.Run(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clone repository"})
	// 	return
	// }
	// gitConfigCmd := exec.Command("git", "config", "--global", "credential.helper", "store")
	// if err := gitConfigCmd.Run(); err != nil {
	// 	// Handle the error
	// 	return
	// }
	// // Write the YAML content to a temporary file
	// tmpfile, err := os.CreateTemp("./tmp/", "knative-service-*.yaml")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file", "details": err.Error()})
	// 	return
	// }
	// if _, err := tmpfile.Write(yamlBytes); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to temporary file"})
	// 	return
	// }

	// // Copy the temporary file to the repository directory

	// // Stage and commit the changes
	// gitAddCmd := exec.Command("git", "add", ".")
	// gitAddCmd.Dir = repoDir
	// if err := gitAddCmd.Run(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stage changes"})
	// 	return
	// }

	// gitCommitCmd := exec.Command("git", "commit", "-m", "Adding knative service YAML")
	// gitCommitCmd.Dir = repoDir
	// if err := gitCommitCmd.Run(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit changes"})
	// 	return
	// }

	// // Set up Git credentials

	// // Push the changes
	// gitPushCmd := exec.Command("git", "push", "--set-upstream", "origin", "main")
	// gitPushCmd.Dir = repoDir
	// var output bytes.Buffer
	// gitPushCmd.Stdout = &output
	// gitPushCmd.Stderr = &output
	// if err := gitPushCmd.Run(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push changes", "details": output.String()})
	// 	return
	// }

	// Set KUBECONFIG environment variable to point to your existing kubeconfig file
	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}
	// cmd := exec.Command("kubectl", "apply", "-f", tmpfile.Name())
	// cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)
	filepath := fmt.Sprintf("./tmp/%s.yaml", fileName)
	outputnew, err := helpers.RetryHelmCommand("kn", []string{"service", "apply", "-f", filepath}, kubeconfigPath, 3, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply resource using kubectl", "details": err.Error()})
		return
	}
	defer os.RemoveAll("./tmp")
	// tenantURL := fmt.Sprintf("https://%s.platform.tripon.io", userInput.Name)
	// toAddress := "envoyfacilitation@gmail.com"
	// subject := "Tenant Deployed"
	// message := fmt.Sprintf("Your tenant %s is deployed successfully", userInput.Name)

	// helpers.CheckTenantUrlSendEmail(tenantURL, toAddress, subject, message, 5*time.Second)

	c.JSON(http.StatusOK, gin.H{"message": " Tenant applied successfully", "details": string(outputnew)})

}

// @BasePath /api/v1
// @Summary get Tenant
// @Schemes
// @Description List Tenant
// @Tags Tenant
// @Accept json
// @Produce json
// @Success 200 {string} gettenant
// @Router /tenantGet [get]
func TenantGet(c *gin.Context) {
	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}

	output, err := helpers.RetryHelmCommand("kn", []string{"service", "list", "-o", "json"}, kubeconfigPath, 3, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply resource using kubectl", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": string(output)})

}

// @Summary Delete Tenant
// @Description Delete Tenant by name
// @Tags Tenant
// @Accept json
// @Produce json
// @Param name body string true "Tenant Name"
// @Success 200 {string} string "Tenant Deleted successfully"
// @Router /tenantDelete [post]
func TenantDelete(c *gin.Context) {
	var userInput struct {
		Name string `json:"name"`
	}
	// implement redis cache

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	kubeconfigPath, err := clients.GetKubeconfigPath(os.Getenv("APP_ENV"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Kubernetes cluster", "details": err.Error()})
		return
	}
	// cmd := exec.Command("kn", "service", "delete", userInput.Name)
	// cmd.Env = append(os.Environ(), "KUBECONFIG="+kubeconfigPath)

	output, err := helpers.RetryHelmCommand("kn", []string{"service", "delete", userInput.Name}, kubeconfigPath, 3, 5*time.Second)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run kn service delete", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": " Tenant Deleted successfully", "details": string(output)})
}
