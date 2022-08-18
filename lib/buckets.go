package lib

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gabriel-vasile/mimetype"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
)

var (
	s3Vultr    *s3.S3
	accessKey  string
	secretKey  string
	bucketName string
	directory  string
	filename   string
	input      int
	input2     string
	loop       string
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
	accessKey = os.Getenv("accessKey")
	secretKey = os.Getenv("secretKey")
	s3Vultr = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("Singapore"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String("https://sgp1.vultrobjects.com/"),
	})))
}

func ListAllBuckets() (ret *s3.ListBucketsOutput) {
	ret, err := s3Vultr.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	return ret
}
func nameCheck(str string) bool {
	checker := regexp.MustCompile(`^[a-z0-9-]*$`).MatchString(str)
	alphanumeric := "abcdefghijklmnopqrstuvwxyz1234567890"
	if checker && strings.Contains(alphanumeric, string(str[0])) && len(str) >= 3 && len(str) <= 63 {
		return true
	} else {
		return false
	}
}

func newBucket() (ret *s3.CreateBucketOutput) {
	_, err := s3Vultr.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(""),
		},
	})
	if awsError, ok := err.(awserr.Error); ok {
		if awsError.Code() == s3.ErrCodeBucketAlreadyExists {
			log.Fatalf("Bucket %q already exists. Error: %v", bucketName, awsError.Code())
		}
	} else {
		log.Printf("Successfully created bucket %q", bucketName)
	}
	return ret
}

func CreateBucket() (ret *s3.CreateBucketOutput) {
	for {
		fmt.Println("Bucket names are unique to their location and must meet the following criteria:")
		fmt.Println("Only lowercase and starts with a letter or number. No spaces.")
		fmt.Println("Bucket name may contain dashes")
		fmt.Println("Must be between 3 and 63 characters long.")
		fmt.Print("Enter your preferred Name for the Bucket: ")
		fmt.Scan(&bucketName)
		if nameCheck(bucketName) {
			break
		} else {
			fmt.Printf("%q does not meet the criteria above. Please try again.", bucketName)
			fmt.Print("\n\n")
			continue
		}
	}
	bucketList := ListAllBuckets()
	if len(bucketList.Buckets) != 0 {
		for _, bucket := range bucketList.Buckets {
			if bucketName == *bucket.Name {
				log.Fatalf("Bucket %q already exists and is owned by you.", bucketName)
			} else {
				newBucket()
				break
			}
		}
	} else {
		newBucket()
	}
	return ret
}

func UploadObject() (ret *s3.PutObjectOutput) {
	fmt.Print("Enter the name of the bucket where you want to upload the File/Object: ")
	fmt.Scan(&bucketName)
	fmt.Print("Enter the Path or Directory where you want to upload the File/Object in the bucket: (e.g., assets/css/): ")
	fmt.Scan(&directory)
	fmt.Print("Enter the Path to the file that you want to upload (e.g., css/styles.css): ")
	fmt.Scan(&filename)

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	var mtype2 string
	if strings.Contains(filename, ".css") {
		mtype2 = "text/css"
	} else {
		mtype, errmime := mimetype.DetectFile(filename)
		if errmime != nil {
			log.Fatalf("Error getting Content-Type: %v", errmime)
		}
		mtype2 = mtype.String()
	}

	log.Println("Uploading Object:", filename)
	ret, err = s3Vultr.PutObject(&s3.PutObjectInput{
		Body:        f,
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path.Join(directory, strings.Split(filename, "/")[strings.Count(filename, "/")])),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(mtype2),
	})

	if err != nil {
		panic(err)
	} else {
		log.Printf("File %q was Uploaded Successfully.", filename)
	}
	return ret
}

func ListObjects() (ret *s3.ListObjectsV2Output) {
	fmt.Print("Enter the name of the bucket to list its file/s: ")
	fmt.Scan(&bucketName)
	ret, err := s3Vultr.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		panic(err)
	}
	return ret
}

func GetObject() {
	fmt.Print("Enter the name of the bucket that contains the file that you want to download: ")
	fmt.Scan(&bucketName)
	fmt.Print("Enter the Path or Directory of the file (e.g. assets/css/): ")
	fmt.Scan(&directory)
	fmt.Print("Enter the name of the file that you want to download(e.g., styles.css): ")
	fmt.Scan(&filename)
	log.Println("Downloading: ", filename)

	ret, err := s3Vultr.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path.Join(directory, strings.Split(filename, "/")[strings.Count(filename, "/")])),
	})

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(ret.Body)
	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
	log.Println("File Downloaded Successfully.")
}

func DeleteObject() (ret *s3.DeleteObjectOutput) {
	fmt.Print("Enter the name of the bucket that contains the file that you want to delete: ")
	fmt.Scan(&bucketName)
	fmt.Print("Enter the Path or Directory of the file in the bucket (e.g., assets/css/): ")
	fmt.Scan(&directory)
	fmt.Print("Enter the name of the file that you want to delete (e.g., styles.css): ")
	fmt.Scan(&filename)
	log.Println("Deleting: ", filename)
	ret, err := s3Vultr.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path.Join(directory, strings.Split(filename, "/")[strings.Count(filename, "/")])),
	})

	if err != nil {
		panic(err)
	} else {
		log.Printf("%q deleted Successfully", filename)
	}
	return ret
}

func DeleteAllObjects() {
	fmt.Scan(&bucketName)
	objectList, errList := s3Vultr.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if errList != nil {
		log.Fatalf("Error Listing objects: %v", errList)
	}
	for _, object := range objectList.Contents {
		log.Printf("Deleting %v", *object.Key)
		_, err := s3Vultr.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(*object.Key),
		})
		log.Printf("%v Deleted Successfully", *object.Key)
		if err != nil {
			log.Fatalf("Error Deleteing Objects.")
		}
	}
	log.Println("All Files deleted successfully.")
}

func DeleteBucket() (ret *s3.DeleteBucketOutput) {
	DeleteAllObjects()
	ret, err := s3Vultr.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			switch awsError.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Fatalf("No Bucket exists with the name '%s'", bucketName)
				log.Println()
			default:
				panic(err)
			}
		}
	}
	log.Printf("Bucket %q deleted Successfully.", bucketName)
	return ret
}
