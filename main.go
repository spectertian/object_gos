package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vultr/govultr"
	"log"
	"object_gos/lib"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	apiKey       = ""
	vc           *govultr.Client
	ctx          = context.Background()
	input        int
	input2       string
	loop         string
	objStorageID int
)

func init() {
	sysType := runtime.GOOS
	path_sour := "/www/web/object_gos/.env"
	if sysType != "linux" {
		path_sour = ".env"
	}
	if err := godotenv.Load(path_sour); err != nil {
		log.Println("No .env file found")
	}
	apiKey = os.Getenv("apiKey")
	vc = govultr.NewClient(nil, apiKey)
}

func listClusters() {
	clusterList, err := vc.ObjectStorage.ListCluster(ctx)
	if err != nil {
		log.Panicf("Error listing clusters: %s", err)
	}
	fmt.Println(clusterList)
	if len(clusterList) <= 1 {
		log.Panic("There are no clusters found to create an Object Storage.")
	}
	log.Printf("List of All Clusters: %+v", clusterList)
}

func CreateObjStorage() {
	var objStorageName string
	clusterID := 2
	fmt.Print("Enter Your Desired Object Storage Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		objStorageName = scanner.Text()
		fmt.Print()
		if objStorageName != "" {
			break
		}
	}
	objStorageNew, errn := vc.ObjectStorage.Create(ctx, clusterID, objStorageName)
	if errn != nil {
		log.Panicf("Error creating storage: %s", errn)
	} else {
		log.Printf("Succesfully created an Object Storage with ID: %s", objStorageNew.ID)
	}
	for {
		objStorageNewInfo, erri := vc.ObjectStorage.Get(ctx, objStorageNew.ID)
		if erri != nil {
			log.Panicf("Error getting Object Storage Information: %s", erri)
		}
		if objStorageNewInfo.Status == "active" {
			log.Print("Your Object Storage is now Active.")
			log.Printf("Object Storage Information: %+v", objStorageNewInfo)
			log.Println("S3 Credentials: ")
			log.Printf("Hostname: %s | Access Key: %s | Secret Key: %s", objStorageNewInfo.S3Hostname, objStorageNewInfo.S3AccessKey, objStorageNewInfo.S3SecretKey)
			break
		}
		log.Printf("The Object Storage is currently %s. Waiting another ten seconds until it becomes active.", objStorageNewInfo.Status)
		time.Sleep(time.Second * 10)
	}
}

func listObjStorage1() {
	objStorageList, err := vc.ObjectStorage.List(ctx, nil)
	log.Printf("The Total Number of Object Storage: %d", len(objStorageList))
	for x := 0; x != len(objStorageList); x++ {
		log.Printf("Object Storage #%d:", x+1)
		log.Printf("List of Object Storage: %+v", objStorageList[x])
		log.Println()
	}
	if err != nil {
		log.Panicf("Error listing Object Storage: %s", err)
	}
}

func listObjStorage2() {
	objStorageList, err := vc.ObjectStorage.List(ctx, nil)
	log.Printf("The Total Number of Object Storage: %d", len(objStorageList))
	for x := 0; x != len(objStorageList); x++ {
		log.Printf("Object Storage #%d:", x+1)
		log.Printf("Date Created: %s", objStorageList[x].DateCreated)
		log.Printf("Object Storage ID: %s", objStorageList[x].ID)
		log.Printf("Object Storage Label: %s", objStorageList[x].Label)
		log.Printf("Object Storage Location: %s", objStorageList[x].Location)
		log.Printf("Object Storage Region: %s", objStorageList[x].RegionID)
		log.Printf("Object Storage Hostname: %s", objStorageList[x].S3Hostname)
		log.Printf("Object Storage Access Key: %s", objStorageList[x].S3AccessKey)
		log.Printf("Object Storage Secret Key: %s", objStorageList[x].S3SecretKey)
		log.Printf("Object Storage Status: %s", objStorageList[x].Status)
		log.Printf("Object Storage Cluster ID: %d", objStorageList[x].ObjectStoreClusterID)
		log.Println()
	}
	if err != nil {
		log.Panicf("Error listing Object Storage: %s", err)
	}
}

func getObjStorage() {
	fmt.Println("To Get your Object Storage's ID, List all of your Object Storage.")
	fmt.Print("Enter Your Object Storage's ID to get its Full Information (e.g. cb676a46-66fd-4dfb-b839-443f2e6c0b60): ")
	fmt.Scan(&objStorageID)
	objStorageGet, err := vc.ObjectStorage.Get(ctx, objStorageID)
	log.Printf("Full information of Object Storage with an ID \"%s\".", objStorageID)
	log.Printf("Object Storage ID: %s", objStorageGet.ID)
	log.Printf("Date Created: %s", objStorageGet.DateCreated)
	log.Printf("Label: %s", objStorageGet.Label)
	log.Printf("Location: %s", objStorageGet.Location)
	log.Printf("Region: %s", objStorageGet.RegionID)
	log.Printf("S3 Hostname: %s", objStorageGet.S3Hostname)
	log.Printf("S3 Access Key: %s", objStorageGet.S3AccessKey)
	log.Printf("S3 Secret Key: %s", objStorageGet.S3SecretKey)
	log.Printf("Status: %s", objStorageGet.Status)
	log.Printf("Cluster ID: %d", objStorageGet.ObjectStoreClusterID)
	log.Println()
	if err != nil {
		log.Panicf("Error Getting Object Storage that has an ID %s: %s", objStorageID, err)
	}
}

func delObjStorage() {
	fmt.Println("To Get your Object Storage's ID, List all of your Object Storage.")
	fmt.Print("Enter the ID of the Object Storage that you want to Delete (e.g. cb676a46-66fd-4dfb-b839-443f2e6c0b60): ")
	fmt.Scan(&objStorageID)
	objStorageDel := vc.ObjectStorage.Delete(ctx, objStorageID)
	if objStorageDel == nil {
		log.Printf("Successfully deleted object storage with an ID \"%s\"", objStorageID)
	}
}

func main() {
	for {
		fmt.Println("Input '1' to List All Buckets in your Object Storage.")
		fmt.Println("Input '2' to Create a new Bucket.")
		fmt.Println("Input '3' to Delete a Bucket.")
		fmt.Println("Input '4' to Upload a File/Object.")
		fmt.Println("Input '5' to List All Files/Objects inside a Bucket.")
		fmt.Println("Input '6' to Download a File/Object from a Bucket.")
		fmt.Println("Input '7' to Delete a File/Object from a Bucket.")
		fmt.Println("Input '8' to Delete All Files/Objects from a Bucket.")
		fmt.Println("Input '9' to Create an Object Storage.")
		fmt.Print("Your Input: ")
		fmt.Scan(&input)
		switch input {
		case 1:
			log.Println(lib.ListAllBuckets())
		case 2:
			lib.CreateBucket()
		case 3:
			fmt.Print("Enter the Name of the Bucket that you want to Delete: ")
			lib.DeleteBucket()
		case 4:
			lib.UploadObject()
		case 5:
			log.Println(lib.ListObjects())
		case 6:
			lib.GetObject()
		case 7:
			lib.DeleteObject()
		case 8:
			fmt.Print("Enter the Name of the Bucket that you want to Empty: ")
			lib.DeleteAllObjects()
		case 9:
			CreateObjStorage()
		default:
			log.Println("Invalid Input! Please try again.")
			continue
		}

		fmt.Print("Do you want to rerun the program? (y/n): ")
		fmt.Scan(&input2)
		loop = strings.ToLower(input2)
		if loop == "n" {
			fmt.Println("Closing the Program...")
			time.Sleep(2 * time.Second)
			break
		} else if loop == "y" {
			continue
		} else {
			log.Fatalln("Invalid Input! Closing the Program.")
			time.Sleep(2 * time.Second)
		}
	}
}
