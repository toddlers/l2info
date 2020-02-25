package lambdalayers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/toddlers/l2info/client"
	"github.com/toddlers/l2info/utils"
)

func ListLayers(runtime string) (*lambda.ListLayersOutput, error) {
	lambdaClient := client.GetClient(nil)
	input := &lambda.ListLayersInput{
		CompatibleRuntime: aws.String(runtime),
	}
	result, err := lambdaClient.ListLayers(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetInfo(layerArn string) (*lambda.GetLayerVersionByArnOutput, error) {
	larn, err := arn.Parse(layerArn)
	if err != nil {
		return nil, err
	}
	sess := client.GetSession(larn.Region)
	lambdaClient := client.GetClient(sess)
	input := &lambda.GetLayerVersionByArnInput{
		Arn: aws.String(larn.String()),
	}
	layerInfo, err := lambdaClient.GetLayerVersionByArn(input)
	if err != nil {
		return nil, err
	}
	return layerInfo, nil
}

func DownloadLayer(exportPath, layerArn string) (string, error) {
	larn, err := arn.Parse(layerArn)
	if err != nil {
		return "", fmt.Errorf("Lambda Layer ARN %s provided is not correct: %v", layerArn, err)
	}
	layerInfo, err := GetInfo(layerArn)
	layerName := strings.Split(larn.Resource, ":")[1]
	if err != nil {
		return "", err
	}
	response, err := http.Get(*layerInfo.Content.Location)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	absPathDir, err := filepath.Abs(exportPath)
	if err != nil {
		return "", fmt.Errorf("Export path %s is not valid : %v", absPathDir, err)
	}
	layerContentZipFile := filepath.Join(absPathDir, layerName+".zip")
	os.MkdirAll(absPathDir, 0755)
	out, err := os.Create(layerContentZipFile)
	if err != nil {
		return "", fmt.Errorf("Can't create layer zip file %s : %v", layerContentZipFile, err)
	}
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", fmt.Errorf("Cant't write to zip file %s : %v", layerContentZipFile, err)
	}
	layerDir := filepath.Join(absPathDir, layerName+"-content")
	err = utils.Unzip(layerContentZipFile, layerDir)
	if err != nil {
		return "", fmt.Errorf("Can't unzip the content %v", err)
	}
	err = os.Remove(layerContentZipFile)
	if err != nil {
		return "", fmt.Errorf("Can't delete zip file %s: %v", layerContentZipFile, err)
	}
	return layerDir, nil
}
