package helpers

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nwaples/rardecode"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	"github.com/tidwall/pretty"
	mail "github.com/xhit/go-simple-mail/v2"
	"gitlab.com/mims-api-service/constants"
	"gitlab.com/mims-api-service/logs"
	"gitlab.com/mims-api-service/models"
	"gitlab.com/mims-api-service/responses"
	"golang.org/x/exp/constraints"
	"gopkg.in/validator.v2"
)

// htmlTemplateCache caches parsed HTML templates by path to avoid repeated ParseFiles and disk I/O.
var htmlTemplateCache = struct {
	sync.RWMutex
	m map[string]*template.Template
}{m: make(map[string]*template.Template)}

// getCachedHtmlTemplate returns a cached template for key, or builds and caches it via build().
func getCachedHtmlTemplate(key string, build func() (*template.Template, error)) (*template.Template, error) {
	htmlTemplateCache.RLock()
	t := htmlTemplateCache.m[key]
	htmlTemplateCache.RUnlock()
	if t != nil {
		return t, nil
	}
	htmlTemplateCache.Lock()
	defer htmlTemplateCache.Unlock()
	if t = htmlTemplateCache.m[key]; t != nil {
		return t, nil
	}
	t, err := build()
	if err != nil {
		return nil, err
	}
	htmlTemplateCache.m[key] = t
	return t, nil
}

func RarExtractor(source string, destination string) error {

	rr, err := rardecode.OpenReader(source, "")

	if err != nil {
		return fmt.Errorf("read: failed to create reader: %v", err)
	}

	for {
		header, err := rr.Next()
		if err == io.EOF {
			break
		}

		if header.IsDir {
			err = Mkdir(filepath.Join(destination, header.Name), 0755)
			if err != nil {
				return err
			}
			continue
		}
		err = Mkdir(filepath.Dir(filepath.Join(destination, header.Name)), 0755)
		if err != nil {
			return err
		}

		err = writeNewFile(filepath.Join(destination, header.Name), rr, header.Mode())
		if err != nil {
			return err
		}
	}
	return nil
}

func Mkdir(path string, dirMode os.FileMode) error {
	err := os.MkdirAll(path, dirMode)
	if err != nil {
		return fmt.Errorf("%s: creating directory: %v", path, err)
	}
	return nil
}

func CheckExt(filename string) error {
	if !strings.HasSuffix(filename, ".rar") {
		return fmt.Errorf("filename must have a .rar extension")
	}
	return nil
}

func writeNewFile(path string, in io.Reader, mode os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("%s: creating directory for file: %v", path, err)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%s: creating new file: %v", path, err)
	}
	defer out.Close()

	err = out.Chmod(mode)
	if err != nil && runtime.GOOS != "windows" {
		return fmt.Errorf("%s: changing file mode: %v", path, err)
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%s: writing file: %v", path, err)
	}
	return nil
}

////////////

// Get filename(s) from within the Archive
func GetRarContents(source string) (string, error) {

	rr, err := rardecode.OpenReader(source, "")

	if err != nil {
		return "", fmt.Errorf("read: failed to create reader: %v", err)
	}

	header, err := rr.Next()
	if err == io.EOF {
		return "", fmt.Errorf("archive is empty: %v", err)
	}
	return header.Name, nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ConverstError(err error) map[string]string {
	stringSlice := strings.Split(err.Error(), ",")
	strErrors := make(map[string]string)
	for _, err := range stringSlice {
		if !strings.Contains(err, ":") {
			continue
		}
		stringSlice2 := strings.Split(err, ":")
		errField := ToSnakeCase(strings.ReplaceAll(stringSlice2[0], " ", ""))
		errMessage := strings.TrimSpace(stringSlice2[1])
		strErrors[errField] = HandleErrMesssage(errMessage)
	}
	return strErrors
}

func HandleErrMesssage(err string) string {
	switch err {
	case "less than min":
		return "โปรดระบุข้อมูลไม่ต่ำกว่า 8 ตัวอักษร"
	case "zero value":
		return "โปรดระบุ"
	case "regular expression mismatch":
		return "รูปแบบไม่ถูกต้อง"
	case "incorrect":
		return "ข้อมูลไม่ถูกต้อง"
	case "mismatch regis":
		return "รหัสผ่านไม่ตรงกัน"
	case "mismatch":
		return "รหัสผ่านไม่ตรงกัน"
	case "mismatch with old password":
		return "รหัสผ่านไม่ตรงกัน"
	case "already used":
		return "ข้อมูลนี้ถูกใช้ในระบบแล้ว"
	case "not found":
		return "ไม่พบข้อมูล"
	case "inactive":
		return "ข้อมูลไม่ได้ใช้งาน"
	case "not verify":
		return "ข้อมูลไม่ได้ยืนยัน"
	case "not equal":
		return "ข้อมูลไม่ครบถ้วน"
	case "duplicate":
		return "ข้อมูลนี้ถูกใช้ในระบบแล้ว"
	case "have space":
		return "โปรดระบุข้อมูลที่ไม่มีช่องว่าง"
	default:
		return "เกิดข้อผิดพลาด โปรดลองอีกครั้ง"
	}
}

func GetSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func IntToInt64(val int) (n int64) {
	n, err := strconv.ParseInt(strconv.Itoa(val), 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	return n
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

func ImplodeInterface(values []interface{}, separator string) string {
	var strValues []string
	for _, v := range values {
		if fmt.Sprintf("%v", reflect.TypeOf(v)) == "string" {
			if fmt.Sprintf("%s", v) != "NULL" {
				strValues = append(strValues, "'"+fmt.Sprintf("%s", v)+"'")
			} else {
				strValues = append(strValues, fmt.Sprintf("%s", v))
			}
		} else if fmt.Sprintf("%v", reflect.TypeOf(v)) == "time.Time" {

			strValues = append(strValues, "'"+v.(time.Time).Format("2006-01-02 15:04:05")+"'")
		} else {
			if fmt.Sprintf("%v", v) == "<nil>" {
				continue
			}
			strValues = append(strValues, fmt.Sprintf("%v", v))
			// fmt.Println(fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(strValues, separator)
}

func ImplodeInterface2(values []interface{}, separator string) string {
	var strValues []string
	for _, v := range values {
		if fmt.Sprintf("%v", reflect.TypeOf(v)) == "string" {
			if fmt.Sprintf("%s", v) != "NULL" {
				strValues = append(strValues, fmt.Sprintf("%s", v))
			} else {
				strValues = append(strValues, fmt.Sprintf("%s", v))
			}
		} else if fmt.Sprintf("%v", reflect.TypeOf(v)) == "time.Time" {

			strValues = append(strValues, "'"+v.(time.Time).Format("2006-01-02 15:04:05")+"'")
		} else {
			if fmt.Sprintf("%v", v) == "<nil>" {
				continue
			}
			strValues = append(strValues, fmt.Sprintf("%v", v))
			// fmt.Println(fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(strValues, separator)
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func PrintlnJson(v ...interface{}) {
	bytes, _ := json.MarshalIndent(v, "", "\t")
	fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))
}

func ConvertStringToInt(input string) (int, error) {
	result, err := strconv.Atoi(input)
	if err != nil {
		return 0, responses.NewAppErr(http.StatusBadRequest, "input params is not a number")
	}
	return result, nil
}

func HasPermission(permissions, accessKeys []string) bool {
	hasPermission := false
	for _, item := range accessKeys {
		if contains(permissions, item) {
			hasPermission = true
			break
		}
	}
	return hasPermission
}

func SortSlice[T constraints.Ordered](list []T, reverse bool) {
	sort.Slice(list, func(i, j int) bool {
		if reverse {
			return list[i] > list[j]
		}
		return list[i] < list[j]
	})
}
func InArray(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func IsFileSizeGreaterThanLimit(fileSize, limitSize int64) bool {
	return fileSize > limitSize*1048576
}

type UserInfo struct {
	UserId         int
	UserDept       int
	UserPermission []string
}

func GetUserInfo(c *gin.Context) UserInfo {
	userInfo, _ := c.MustGet("userInfo").(UserInfo)
	return userInfo
}

func GetUserPermission(payload []interface{}) []string {
	permissionSlice := []string{}

	for _, permission := range payload {
		permissionSlice = append(permissionSlice, permission.(string))
	}

	return permissionSlice
}

type EmailContent struct {
	SendFromEmail string
	SendToEmail   string
	EmailSubject  string
	EmailMessage  string
}

func CreateEmailAndSend(sendEmail EmailContent) error {
	// Connect to server
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_HOST")
	server.Port = port
	server.Username = os.Getenv("SMTP_USER")
	server.Password = os.Getenv("SMTP_PASS")
	server.Encryption = mail.EncryptionSTARTTLS
	server.KeepAlive = false

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Create email
	emailMsg := mail.NewMSG()
	emailMsg.SetFrom(sendEmail.SendFromEmail)
	emailMsg.AddTo(sendEmail.SendToEmail)
	emailMsg.SetSubject(sendEmail.EmailSubject)
	emailMsg.SetBody(mail.TextHTML, sendEmail.EmailMessage)

	// emailMsg.Attach(&mail.File{FilePath: "./assets/logo.png", Name: "logo.png", Inline: true})

	// Send email
	sendErr := emailMsg.Send(smtpClient)
	if sendErr != nil {
		return sendErr
	}

	return nil
}

func GetUniqueKey(keys map[float32]bool) []float32 {
	keysSlice := []float32{}
	for key := range keys {
		keysSlice = append(keysSlice, key)
	}

	sort.SliceStable(keysSlice, func(i, j int) bool {
		return int(keysSlice[i]) < int(keysSlice[j])
	})

	return keysSlice
}

func DamageSortKm(keys []responses.RoadDamageDetail, direction int) []responses.RoadDamageDetail {
	if direction == 1 {
		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i].KmStart < keys[j].KmStart
		})
	} else {
		sort.SliceStable(keys, func(i, j int) bool {
			return keys[i].KmStart > keys[j].KmStart
		})
	}

	return keys
}
func Min(value float64, values ...float64) float64 {
	for _, v := range values {
		if v < value {
			value = v
		}
	}
	return value
}

func Max(value float64, values ...float64) float64 {
	for _, v := range values {
		if v > value {
			value = v
		}
	}
	return value
}

func RoundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func RoundFloatPointer(val *float64, precision int) {
	if val == nil {
		return
	}
	*val = math.Round(*val*math.Pow10(precision)) / math.Pow10(precision)
}

func RoundStructFloats(s interface{}, precision int) reflect.Value {
	inputValue := reflect.ValueOf(s)
	inputType := inputValue.Type()
	structValue := inputValue

	output := reflect.New(inputType).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)

		if field.Kind() == reflect.Float64 {
			roundedValue := RoundFloat(field.Float(), precision)
			output.Field(i).SetFloat(roundedValue)
		} else {
			output.Field(i).Set(field)
		}
	}

	return output
}
func RoundStructPointerFloats(s interface{}, precision int) interface{} {
	inputValue := reflect.ValueOf(s)
	inputType := inputValue.Type()

	if inputValue.Kind() != reflect.Ptr {
		return nil
	}

	structValue := inputValue.Elem()
	output := reflect.New(inputType.Elem()).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)

		if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Float64 {
			if field.IsNil() {
				continue
			}
			roundedValue := RoundFloat(field.Elem().Float(), precision)
			roundedField := reflect.New(field.Type().Elem())
			roundedField.Elem().SetFloat(roundedValue)
			output.Field(i).Set(roundedField)
		} else {
			output.Field(i).Set(field)
		}
	}

	return output.Addr().Interface()

}

func SaveFile(c *gin.Context, file *multipart.FileHeader, dir string) (string, error) {
	//add time to make image file unique
	dstPath := dir + time.Now().Format("20060102150405") + "_" + file.Filename
	err := c.SaveUploadedFile(file, dstPath)
	if err != nil {
		return "", err
	}

	return dstPath, nil
}

func StrToFloat(str string) float64 {
	float, _ := strconv.ParseFloat(str, 32)
	return float
}

func StrToFloatValidate(str string) float64 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		//ให้ค่าติดเพื่อเช็คบรรทัดตอน validate กรณีมีข้อมูลเป็นตัวหนังสือ
		value = -1
	}
	return value
}

func StrToFloatPointerValidate(str string) *float64 {
	if str == "" {
		return nil
	}
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		//ให้ค่าติดเพื่อเช็คบรรทัดตอน validate กรณีมีข้อมูลเป็นตัวหนังสือ
		value = -2
	}
	return &value
}

func StrToInt(str string) int {
	intData, _ := strconv.Atoi(str)
	return intData
}

func CheckInt(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func CheckFloat(str string) bool {
	_, err := strconv.ParseFloat(str, 32)
	return err == nil
}

func Unzip(dstPath, dir string) error {

	fileType := GetFileType(dstPath)
	if fileType == ".rar" {
		err := RarExtractor(dstPath, dir)
		if err != nil {
			return err
		}
	} else {

		// output := "storages/road/damage/test"
		archive, err := zip.OpenReader(dstPath)
		if err != nil {
			return err
		}

		defer archive.Close()
		for _, f := range archive.File {
			filePath := filepath.Join(dir, f.Name)
			if strings.Contains(filePath, "__MACOSX") {
				continue
			}

			if !strings.HasPrefix(filePath, filepath.Clean(dir)+string(os.PathSeparator)) {
				fmt.Println("invalid file path")
				return errors.New("invalid file path")
			}
			if f.FileInfo().IsDir() {
				fmt.Println("creating directory...", filePath)
				os.MkdirAll(filePath, 0775)
				continue
			}

			// if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			// 	return err
			// }

			dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0775)
			if err != nil {
				return err
			}

			// buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
			// // _, err = dstFile.Read(buff)
			// // if err != nil {
			// // 	fmt.Println("ake([]byte, 512)", err)
			// // 	os.Exit(1)
			// // }
			// // filetype := http.DetectContentType(buff)

			// fmt.Println(filetype)

			fileInArchive, err := f.Open()
			if err != nil {
				return err
			}

			if _, err := io.Copy(dstFile, fileInArchive); err != nil {
				return err
			}

			dstFile.Close()
			fileInArchive.Close()
		}
	}
	return nil
}

func GetFileType(filePath string) string {
	extension := filepath.Ext(filePath)
	return extension
}

func ContainsInt(item int, items []int) bool {
	for _, value := range items {
		if item == value {
			return true
		}
	}
	return false
}

func TypeFileAllowed(path, allowedTypes string) (bool, error) {
	data := Explode(".", path)
	if len(data) <= 1 {
		return false, errors.New("file Content Type not found")
	}
	types := Explode("|", allowedTypes)
	checkType := stringInSlice(data[1], types)
	if !checkType {
		return false, errors.New("the filetype you are attempting to upload is not allowed")
	}
	return checkType, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func SortStructByID(list interface{}, reverse bool) error {
	if list == nil {
		return nil
	}

	// Get the type and value of the input list
	listValue := reflect.ValueOf(list)
	listType := listValue.Type()

	if listType.Kind() != reflect.Slice {
		return errors.New("SortByID: non-slice input")
	}

	// Check if the element type of the list has an ID field
	elemType := listType.Elem()
	_, hasIDField := elemType.FieldByName("ID")
	if !hasIDField {
		return errors.New("SortByID: struct does not have an ID field")

	}

	// Define the sorting function
	less := func(i, j int) bool {
		iID := listValue.Index(i).FieldByName("ID").Int()
		jID := listValue.Index(j).FieldByName("ID").Int()
		if reverse {
			return iID > jID
		}
		return iID < jID
	}

	// Use reflection to sort the list
	sort.Slice(list, less)
	return nil
}

func MoveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFile(sourcePath string) error {
	err := os.Remove(sourcePath)
	if err != nil {
		return err
	}
	return nil
}

func GetAccessControl(c *gin.Context) []string {
	permissions := []string{}
	userPermission, _ := c.Get("accessControl")
	if userPermission != nil {
		for _, item := range userPermission.([]interface{}) {
			permissions = append(permissions, item.(string))
		}
	}
	return permissions
}

func GetTimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func SortedMap(m map[string]interface{}) map[string]interface{} {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// create a new map with the sorted key-value pairs
	sortedMap := make(map[string]interface{})
	for _, k := range keys {
		sortedMap[k] = m[k]
	}
	return sortedMap
}

func DownloadHandler(c *gin.Context, pathFile string) {
	file, err := os.Open(pathFile)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the appropriate headers to tell the client that it should download the file as a ZIP archive
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=file.zip")

	// Serve the file to the client using the File() function
	c.File(file.Name())
}

func DecodeImgBase64(imgBase64, filePath, fileName, fileType string) (string, error) {
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		panic(err)
	}
	// Decode base64-encoded string
	imgData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(imgBase64))
	if err != nil {
		return "", err
	}
	if fileType == "png" {
		// Decode image data into image.Image
		img, fileType, err := image.Decode(bytes.NewReader(imgData))
		if err != nil {
			return "", err
		}

		// Save the image to a file
		// uuid := uuid.New()
		out, err := os.Create(fmt.Sprintf("%s.%s", filePath+fileName, fileType))
		if err != nil {
			return "", errors.New(constants.FAILED_TO_SAVE_FILE)
		}
		defer out.Close()

		// Encode the image as PNG and write it to the output file
		err = png.Encode(out, img)
		if err != nil {
			return "", err
		}
	} else {
		outputFile, err := os.Create(filePath + fileName + "." + fileType)
		if err != nil {
			return "", err
		}
		defer outputFile.Close()

		// Write the decoded content to the output file
		_, err = outputFile.Write(imgData)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%s.%s", filePath+fileName, fileType), nil
}

func DecodeFileBase64FileType(imgBase64, filePath, fileName string, fileType string, file string) (string, string, error) {
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		panic(err)
	}
	if file == "file" {
		if fileType != "docx" && fileType != "xlsx" && fileType != "pdf" && fileType != "dwg" && fileType != "jpeg" && fileType != "jpg" && fileType != "png" {
			return "", "", errors.New("ivalid file type. Only docx, xlsx, pdf, dws, jpeg, jpg and png files are allowed")
		}
	}

	if file == "img" {
		if fileType != "jpeg" && fileType != "jpg" && fileType != "png" {
			return "", "", errors.New("ivalid file type. Only jpeg, jpg and png files are allowed")
		}
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(imgBase64))
	if err != nil {
		logs.Error(err)
		return "", "", err
	}
	switch fileType {
	case "docx":
		err = ioutil.WriteFile(filePath+fileName+".docx", data, 0777)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
	case "xlsx":
		err = ioutil.WriteFile(filePath+fileName+".xlsx", data, 0777)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
	case "pdf":
		err = ioutil.WriteFile(filePath+fileName+".pdf", data, 0777)
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
	case "dwg":
		outputFile, err := os.Create(filePath + fileName + ".dwg")
		if err != nil {
			return "", "", err
		}
		defer outputFile.Close()

		// Write the decoded content to the output file
		_, err = outputFile.Write(data)
		if err != nil {
			return "", "", err
		}
	case "jpeg":
		outputFile, err := os.Create(filePath + fileName + ".jpeg")
		if err != nil {
			return "", "", err
		}
		defer outputFile.Close()

		// Write the decoded content to the output file
		_, err = outputFile.Write(data)
		if err != nil {
			return "", "", err
		}
	case "jpg":
		outputFile, err := os.Create(filePath + fileName + ".jpg")
		if err != nil {
			return "", "", err
		}
		defer outputFile.Close()

		// Write the decoded content to the output file
		_, err = outputFile.Write(data)
		if err != nil {
			return "", "", err
		}
	case "png":
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			logs.Error(err)
			return "", "", err
		}

		out, err := os.Create(fmt.Sprintf("%s.%s", filePath+fileName, fileType))
		if err != nil {
			logs.Error(err)
			return "", "", errors.New(constants.FAILED_TO_SAVE_FILE)
		}
		defer out.Close()

		err = png.Encode(out, img)
		if err != nil {
			logs.Error(err)
			return "", "", err
		}
	}

	return fmt.Sprintf("%s.%s", filePath+fileName, fileType), fileType, nil
}

func DecodeImgBase64FileType(imgBase64, filePath, fileName string, typeFile string) (string, string, error) {
	err := os.MkdirAll(filePath, 0777)
	if err != nil {
		panic(err)
	}
	// PrintlnJson(fileName)
	// if typeFile == "pdf" {
	// 	err := decodeBase64ToPDF(imgBase64, filePath+fileName)
	// 	return "", "", err
	// }
	// return "", "", nil
	// Decode base64-encoded string
	imgData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(imgBase64))
	if err != nil {
		logs.Error(err)
		return "", "", err
	}
	PrintlnJson(imgData)
	err = ioutil.WriteFile(filePath+fileName+".pdf", imgData, 0777)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	// Decode image data into image.Image
	img, fileType, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		logs.Error(err)
		return "", "", err
	}
	if typeFile == "file" {
		if fileType != "docx" && fileType != "xlsx" && fileType != "pdf" && fileType != "dws" && fileType != "jpeg" && fileType != "jpg" && fileType != "png" {
			return "", "", errors.New("ivalid file type. Only docx, xlsx, pdf, dws, jpeg, jpg and png files are allowed")
		}
	}

	if typeFile == "img" {
		if fileType != "jpeg" && fileType != "jpg" && fileType != "png" {
			return "", "", errors.New("ivalid file type. Only jpeg, jpg and png files are allowed")
		}
	}

	// Save the image to a file
	// uuid := uuid.New()
	out, err := os.Create(fmt.Sprintf("%s.%s", filePath+fileName, fileType))
	if err != nil {
		logs.Error(err)
		return "", "", errors.New(constants.FAILED_TO_SAVE_FILE)
	}
	defer out.Close()

	// Encode the image as PNG and write it to the output file
	err = png.Encode(out, img)
	if err != nil {
		logs.Error(err)
		return "", "", err
	}
	return fmt.Sprintf("%s.%s", filePath+fileName, fileType), fileType, nil
}

func RemoveEmptySlice(slice interface{}, fieldName string) (interface{}, error) {
	s := reflect.ValueOf(slice)

	if s.Kind() != reflect.Slice {
		return nil, errors.New("removeEmptySlice() given a non-slice type")
	}

	var newSlice reflect.Value

	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i)
		subSlice := elem.FieldByName(fieldName)

		// Check if the field exists
		if !subSlice.IsValid() {
			// Return an error if the field does not exist
			return nil, errors.New("the provided field does not exist")
		}

		if subSlice.Len() > 0 {
			if !newSlice.IsValid() {
				newSlice = reflect.MakeSlice(s.Type(), 0, s.Len())
			}
			newSlice = reflect.Append(newSlice, elem)
		}
	}

	if !newSlice.IsValid() {
		return reflect.MakeSlice(s.Type(), 0, 0).Interface(), nil
	}
	return newSlice.Interface(), nil
}

func StringToGeom(str string) (string, error) {
	wkbBytes, err := hex.DecodeString(str)
	if err != nil {
		return "", err
	}
	geom, err := wkb.Unmarshal(wkbBytes)
	if err != nil {
		return "", err
	}
	wktString := wkt.MarshalString(geom)
	return wktString, nil
}

func HasNaN(response interface{}) bool {
	value := reflect.ValueOf(response)

	switch value.Kind() {
	case reflect.Ptr:
		return HasNaN(value.Elem().Interface())
	case reflect.Slice, reflect.Array:
		for i := 0; i < value.Len(); i++ {
			if HasNaN(value.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Struct:
		for i := 0; i < value.NumField(); i++ {
			if HasNaN(value.Field(i).Interface()) {
				return true
			}
		}
	case reflect.Float32, reflect.Float64:
		if math.IsNaN(value.Float()) {
			return true
		}
	}
	return false
}

func SortStructByField(list interface{}, fieldName string, reverse bool) error {
	if list == nil {
		return nil
	}

	// Get the type and value of the input list
	listValue := reflect.ValueOf(list)
	listType := listValue.Type()

	if listType.Kind() != reflect.Slice {
		return errors.New("SortByField: non-slice input")
	}

	// Check if the element type of the list has the specified field
	elemType := listType.Elem()
	_, hasField := elemType.FieldByName(fieldName)
	if !hasField {
		return errors.New("SortByField: struct does not have the specified field")
	}

	// Define the sorting function
	less := func(i, j int) bool {
		iVal := listValue.Index(i).FieldByName(fieldName).Interface()
		jVal := listValue.Index(j).FieldByName(fieldName).Interface()

		// Compare numeric values
		if iNum, ok := iVal.(int); ok {
			jNum, _ := jVal.(int)
			if reverse {
				return iNum > jNum
			}
			return iNum < jNum
		}

		// Compare string values
		if iStr, ok := iVal.(string); ok {
			jStr, _ := jVal.(string)
			if reverse {
				return iStr > jStr
			}
			return iStr < jStr
		}

		// Compare other types
		iStr := fmt.Sprintf("%v", iVal)
		jStr := fmt.Sprintf("%v", jVal)
		if reverse {
			return iStr > jStr
		}
		return iStr < jStr
	}

	// Use reflection to sort the list
	sort.Slice(list, less)
	return nil
}

func CountDecimal(f float64) int {
	s := strconv.FormatFloat(f, 'g', -1, 64)

	// Find the position of the decimal point.
	point := strings.Index(s, ".")

	// If there's no decimal point, there are no decimal places.
	if point == -1 {
		return 0
	}

	// Count the characters after the decimal point, which are the decimal places.
	return len(s) - point - 1
}

func SortStructByFieldFloat32(list interface{}, fieldName string, reverse bool) error {
	if list == nil {
		return nil
	}

	// Get the type and value of the input list
	listValue := reflect.ValueOf(list)
	listType := listValue.Type()

	if listType.Kind() != reflect.Slice {
		return errors.New("SortByField: non-slice input")
	}

	// Check if the element type of the list has the specified field
	elemType := listType.Elem()
	_, hasField := elemType.FieldByName(fieldName)
	if !hasField {
		return errors.New("SortByField: struct does not have the specified field")
	}

	// Define the sorting function
	less := func(i, j int) bool {
		iVal := listValue.Index(i).FieldByName(fieldName).Interface()
		jVal := listValue.Index(j).FieldByName(fieldName).Interface()

		// Compare numeric values
		if iNum, ok := iVal.(float32); ok {
			jNum, _ := jVal.(float32)
			if reverse {
				return iNum > jNum
			} else {
				return iNum < jNum
			}
		}

		// Compare string values
		if iStr, ok := iVal.(string); ok {
			jStr, _ := jVal.(string)
			if reverse {
				return iStr > jStr
			} else {
				return iStr < jStr
			}
		}

		// Compare other types
		iStr := fmt.Sprintf("%v", iVal)
		jStr := fmt.Sprintf("%v", jVal)
		if reverse {
			return iStr > jStr
		} else {
			return iStr < jStr
		}
	}

	// Use reflection to sort the list
	sort.Slice(list, less)
	return nil
}

func SetPermissionsTo0775(path string) error {
	err := os.Chmod(path, 0775)
	return err
}

func FindMinMaxFloat64(values []float64) (float64, float64) {
	min := values[0]
	max := values[0]
	for _, number := range values {
		if number < min {
			min = number
		}
		if number > max {
			max = number
		}
	}
	return min, max
}

func FindMaxInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, number := range values {
		if number > max {
			max = number
		}
	}
	return max
}

func SetTimeToString(t time.Time) string {
	var result string
	if t.IsZero() {
		result = ""

	} else {
		result = t.Format(time.RFC3339)
	}
	return result
}

func CalculateRoadConditionAverageAndCount(dvCount *float64, target *float64, source *float64, kmStart float64, kmEnd float64) (*float64, *float64) {

	if source == nil {
		return nil, nil
	}
	if target == nil {
		if *source >= 0 {
			val := 0.0
			target = &val
		} else {
			val := 0.0
			return &val, nil
		}
	}
	if dvCount == nil {
		initialValue := 0.0
		dvCount = &initialValue
	}

	*dvCount += math.Abs(kmStart - kmEnd)

	*target += *source * math.Abs(kmStart-kmEnd)
	return target, dvCount

}

func CalculateRoadConditionAverage(target *float64, source *float64, kmStart float64, kmEnd float64) *float64 {

	if source == nil {
		return nil
	}
	if target == nil {
		if *source >= 0 {
			val := 0.0
			target = &val
		} else {
			val := 0.0
			return &val
		}
	}

	*target += *source * math.Abs(kmStart-kmEnd)
	return target

}

func CalculateRoadCondition(numerator, denominator *float64) *float64 {
	if numerator != nil && denominator != nil && *denominator > 0 {
		result := *numerator / *denominator
		return &result
	}
	return nil
}

func ConvertNullableFloat64(value *float64) string {
	if value != nil {
		return fmt.Sprintf("%f", *value)
	}
	return "NULL"
}

func CheckFloatIntNonNil(request interface{}, skipFields ...string) error {
	errMsg := ""
	value := reflect.ValueOf(request)
	for i := 0; i < value.NumField(); i++ {
		// Get the current field name
		fieldName := value.Type().Field(i).Name
		// Check if current field is in the list of fields to be skipped
		if stringInSlice(fieldName, skipFields) {
			continue
		}
		field := value.Field(i)
		if field.Type() == reflect.TypeOf((*float64)(nil)) || field.Type() == reflect.TypeOf((*int)(nil)) {
			if field.IsNil() {
				msg := fmt.Sprintf("%s: zero value,", fieldName)
				errMsg += msg
			}
		}
	}
	if errMsg != "" {
		if validateErr := validator.Validate(request); validateErr != nil {
			errMsg2 := validateErr.Error()
			errMsg += errMsg2
			return errors.New(errMsg)
		}
		errMsg = errMsg[:len(errMsg)-1]
		return errors.New(errMsg)
	}
	return nil
}

func CopyPointerToValueFloatInt(toValue interface{}, fromValue interface{}) {
	to := reflect.ValueOf(toValue).Elem()
	from := reflect.ValueOf(fromValue).Elem()

	for i := 0; i < to.NumField(); i++ {
		for j := 0; j < from.NumField(); j++ {
			if to.Type().Field(i).Name == from.Type().Field(j).Name {
				if to.Field(i).Kind() == reflect.Float64 && from.Field(j).Kind() == reflect.Ptr && from.Field(j).Elem().Kind() == reflect.Float64 {
					to.Field(i).SetFloat(from.Field(j).Elem().Float())
				} else if to.Field(i).Kind() == reflect.Int && from.Field(j).Kind() == reflect.Ptr && from.Field(j).Elem().Kind() == reflect.Int {
					to.Field(i).SetInt(from.Field(j).Elem().Int())
				}
			}
		}
	}

	// return
}

func StringToDate(date string) (time.Time, error) {
	layout := "2006-01-02" // The layout pattern corresponding to the date string format
	// Parse the string into a time.Time value
	data, err := time.Parse(layout, date)
	if err != nil {
		// fmt.Println("Error parsing date:", err)
		return data, err
	}
	return data, nil
}

func StringToDateTime(date string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	data, err := time.Parse(layout, date)
	if err != nil {
		return data, err
	}
	return data, nil
}

func InitDataToHtml(templateName string, data interface{}, filePath string) (string, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"minus": func(a, b int) int {
			return a - b
		},
		"formatFloat":         FormatFloat,
		"calculateZoom":       CalculateZoom,
		"formatFloatSetPoint": FormatFloatSetPoint,
	}
	templateGen, err := getCachedHtmlTemplate(templateName, func() (*template.Template, error) {
		return template.New(filepath.Base(templateName)).Funcs(funcMap).ParseFiles(templateName)
	})
	if err != nil {
		PrintlnJson(err.Error())
		logs.Error(err)
		return "", err
	}
	filePath = filePath + "temp/"
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		PrintlnJson(err.Error())
		logs.Error(err)
		return "", err
	}

	fileName := filePath + uuid.New().String() + ".html"
	fileWritter, err := os.Create(fileName)
	if err != nil {
		PrintlnJson(err.Error())
		logs.Error(err)
		return "", err
	}

	if err := templateGen.Execute(fileWritter, data); err != nil {
		PrintlnJson(err.Error())
		logs.Error(err)
		fmt.Println(err)
		return "", err
	}
	if err := fileWritter.Close(); err != nil {
		return "", err
	}
	return fileName, nil
}

func StackValue(data map[string]string, val string) map[string]string {
	data[val] = val
	return data
}

func InitDataToHtmlFunc(templateName string, data interface{}, filePath, filename string) (string, error) {
	var tmplFuncMap = template.FuncMap{
		"convertYear": func(num int) int {
			return num + 543
		},
		"increment": func(num int) int {
			return num + 1
		},
		"km": func(num int) string {
			n := int64(num)
			if n < 0 {
				return "-" + FormatKM(-n)
			} else if n < 1000 {
				in := []byte(strconv.FormatInt(n, 10))
				if len(in) == 1 {
					return fmt.Sprintf("0+00%d", n)
				} else if len(in) == 2 {
					return fmt.Sprintf("0+0%d", n)
				} else {
					return fmt.Sprintf("0+%d", n)
				}
			} else {
				in := []byte(strconv.FormatInt(n, 10))
				var out []byte
				if i := len(in) % 3; i != 0 {
					if out, in = append(out, in[:i]...), in[i:]; len(in) > 0 {
						out = append(out, '+')
					}
				}
				for len(in) > 0 {
					if out, in = append(out, in[:3]...), in[3:]; len(in) > 0 {
						out = append(out, '+')
					}
				}
				return string(out)
			}
		},
	}

	// Create a new template with the custom template function
	// templateGen := template.New("html").

	// Parse the HTML template
	// templateGen, err := templateGen.ParseFiles(templateName)
	// if err != nil {
	// 	return "", err
	// }

	// templateGen, err := template.ParseFiles(templateName).Funcs(tmplFuncMap)
	// if err != nil {
	// 	return "", err
	// }

	htmlTmpl, err := getCachedHtmlTemplate(templateName+"|"+filename, func() (*template.Template, error) {
		return template.New(filename).Funcs(tmplFuncMap).ParseFiles(templateName)
	})
	if err != nil {
		panic(err)
	}

	filePath = filePath + "temp/"
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		logs.Error(err)
		return "", err
	}

	fileName := filePath + uuid.New().String() + ".html"
	fileWritter, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	if err := htmlTmpl.Execute(fileWritter, data); err != nil {
		fmt.Println(err)
		return "", err
	}
	if err := fileWritter.Close(); err != nil {
		PrintlnJson(err.Error())
		logs.Error(err)
		return "", err
	}
	return fileName, nil
}

func PrintToPDF(html string, res *[]byte, isDelay bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}
			delay := 5
			if isDelay {
				delay = 10
			}

			defer chromedp.Run(
				ctx,
				RunWithTimeOut(&ctx, time.Duration(delay), chromedp.Tasks{
					chromedp.WaitReady("div#success-pagejs"),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithDisplayHeaderFooter(true).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}

// NewChromedpContext creates a chromedp context. If CHROME_PATH env is set (e.g. in Docker with Chromium), uses it so PDF generation works.
func NewChromedpContext(parent context.Context, logf func(string, ...interface{})) (context.Context, context.CancelFunc) {
	if path := os.Getenv("CHROME_PATH"); path != "" {
		opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ExecPath(path))
		allocCtx, allocCancel := chromedp.NewExecAllocator(parent, opts...)
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(logf))
		return ctx, func() { cancel(); allocCancel() }
	}
	return chromedp.NewContext(parent, chromedp.WithLogf(logf))
}

func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

func IsString(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.String
}

func SliceIntToString(intSlice []int) string {
	var strSlice []string
	for _, num := range intSlice {
		strSlice = append(strSlice, strconv.Itoa(num))
	}
	return strings.Join(strSlice, ",")
}

func InArrayInt(value int, array []int) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}

func FormatKM(n int64) string {
	if n < 0 {
		return "-" + FormatKM(-n)
	} else if n < 1000 {
		in := []byte(strconv.FormatInt(n, 10))
		if len(in) == 1 {
			return fmt.Sprintf("0+00%d", n)
		} else if len(in) == 2 {
			return fmt.Sprintf("0+0%d", n)
		} else {
			return fmt.Sprintf("0+%d", n)
		}
	} else {
		in := []byte(strconv.FormatInt(n, 10))
		var out []byte
		if i := len(in) % 3; i != 0 {
			if out, in = append(out, in[:i]...), in[i:]; len(in) > 0 {
				out = append(out, '+')
			}
		}
		for len(in) > 0 {
			if out, in = append(out, in[:3]...), in[3:]; len(in) > 0 {
				out = append(out, '+')
			}
		}
		return string(out)
	}
}

func FindValueInArrStr(data []string, targetValue string) ([]string, int) {

	arr := []string{}
	cnt := 0

	for _, value := range data {
		if value == targetValue {
			cnt++
		} else {
			arr = append(arr, value)
		}
	}
	return arr, cnt
}

func GetFileTypeBase64(base64Data string) (string, error) {
	response, err := http.Get(base64Data)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	contentType := response.Header.Get("Content-Type")
	fileType := strings.Split(contentType, "/")[1]

	return fileType, nil
}

func CountDecimalPlaces(number float64) int {
	s := strconv.FormatFloat(number, 'f', -1, 64)
	i := strings.IndexByte(s, '.')
	if i > -1 {
		return len(s) - i - 1
	}
	return 0
}

func GeoJson(geoJSON []byte, index int, data []models.RoadConditionSurveyM2) (responses.Geometry, error) {
	var geometry responses.Geometry

	err := json.Unmarshal(data[index].Geojson, &geometry)
	if err != nil {
		fmt.Println("Error:", err)
		return geometry, err
	}
	return geometry, nil
}

func ConvertToArrayInt(arrayString []string) ([]int, error) {
	intSlice := make([]int, 0, len(arrayString))
	for _, valueStr := range arrayString {
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return nil, err // Conversion failed
		}
		intSlice = append(intSlice, value)
	}
	return intSlice, nil
}

func Float64Ptr(value float64) *float64 {
	return &value
}

func CheckOriginDestinationRoad(refRoadTypeId int, nameOriginTH, nameDestinationTH, roadName string) string {
	var roadOrigin string
	if refRoadTypeId == 1 || refRoadTypeId == 3 {
		if nameDestinationTH != "" {
			roadOrigin = nameOriginTH + " - " + nameDestinationTH
		} else {
			roadOrigin = nameOriginTH
		}
	} else if refRoadTypeId == 2 || refRoadTypeId == 4 {
		if nameDestinationTH != "" {
			roadOrigin = nameDestinationTH + " - " + nameOriginTH
		} else {
			roadOrigin = nameOriginTH
		}
	} else {
		roadOrigin = roadName
	}
	return roadOrigin
}

func GetUserID(c *gin.Context) int {
	userID, _ := c.MustGet("userID").(float64)
	return int(userID)
}

func IsValidDate(date string) bool {
	// Assuming the desired date format is "YYYY-MM-DD"
	_, err := time.Parse("2006-01-02 15:04:05", date)
	return err == nil
}

// func EnsureDir(dirName string) error {
// 	// Check if the directory exists
// 	info, err := os.Stat(dirName)

// 	// If the directory does not exist, create it
// 	if os.IsNotExist(err) {
// 		return os.MkdirAll(dirName, 0775) // You can adjust the file permissions as needed
// 	}

// 	// If the path exists but is not a directory, return an error
// 	if info != nil && !info.IsDir() {
// 		return fmt.Errorf("path exists but is not a directory: %s", dirName)
// 	}

// 	return nil
// }

func KmStringToFloat64(s string) (float64, error) {
	cleanedString := strings.ReplaceAll(s, "+", "")
	return strconv.ParseFloat(cleanedString, 64)
}

func Float64ToKmString(f float64) string {
	// Convert the float64 to an integer string with no decimal places
	str := strconv.FormatFloat(f, 'f', 0, 64) // Force no decimal part

	// Start inserting "+" from the right, every three digits
	var result strings.Builder
	for i, count := len(str), 0; i > 0; i, count = i-1, count+1 {
		// Every three characters, prepend a "+" unless it's the very beginning of the string
		if count%3 == 0 && count > 0 {
			result.WriteString("+")
		}
		// Prepend the current character from str
		result.WriteByte(str[i-1])
	}

	// Reverse the string in the result to correct the order
	finalStr := reverseStrings(result.String())
	return finalStr
}

func reverseStrings(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func FilterByPrefix(slice []string, prefix string) []string {
	filteredSlice := make([]string, 0, len(slice))

	for _, item := range slice {
		if !strings.HasPrefix(item, prefix) {
			filteredSlice = append(filteredSlice, item)
		}
	}

	return filteredSlice
}

func HasDuplicate(arr1, arr2 []int) bool {
	set := make(map[int]bool)

	// Add all elements of the first array to the set
	for _, num := range arr1 {
		set[num] = true
	}
	// Check if any element in the second array is already in the set
	for _, num := range arr2 {
		if set[num] {
			return true
		}
	}

	return false
}

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return &os.PathError{Op: "copy", Path: src, Err: os.ErrInvalid}
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func EnsureDir(dirName string) error {
	// Check if the directory exists
	info, err := os.Stat(dirName)

	// If the directory does not exist, create it
	if os.IsNotExist(err) {
		return os.MkdirAll(dirName, 0775) // You can adjust the file permissions as needed
	}

	// If the path exists but is not a directory, return an error
	if info != nil && !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", dirName)
	}

	return nil
}

func RemoveFileIfExists(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		// The file exists, attempt to remove it.
		if err := os.Remove(filePath); err != nil {
			return err
		}
	}
	return nil
}

func FindMin(nums []int) int {
	if len(nums) == 0 {
		return math.MaxInt
	}

	minVal := nums[0]
	for _, num := range nums {
		if num < minVal {
			minVal = num
		}
	}
	return minVal
}

func ConvertThegeomToGeomJSON(theGeom []byte) (responses.GeomJSON, error) {
	var geometry responses.GeomJSON
	err := json.Unmarshal(theGeom, &geometry)
	if err != nil {
		return responses.GeomJSON{}, err
	}
	return geometry, nil
}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// PrintToPDFWithDelay renders HTML and prints to PDF. readySelector: e.g. "div#success" or "div#success-pagejs"; if empty, defaults to "div#success". delaySeconds: max wait for ready.
func PrintToPDFWithDelay(html string, res *[]byte, delaySeconds int, readySelector string) chromedp.Tasks {
	if readySelector == "" {
		readySelector = "div#success"
	}
	sel := readySelector
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {

			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}

			defer chromedp.Run(
				ctx,
				RunWithTimeOut(&ctx, time.Duration(delaySeconds), chromedp.Tasks{
					chromedp.WaitVisible(sel),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithDisplayHeaderFooter(true).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
func IsContainsSlice[T constraints.Ordered](item T, items []T) bool {
	for _, value := range items {
		if item == value {
			return true
		}
	}
	return false
}

func CalculateZoom(numberOfData int) int {
	switch {
	case numberOfData > 600000:
		return 20
	case numberOfData > 500000:
		return 15
	case numberOfData > 300000:
		return 12
	case numberOfData > 100000:
		return 11
	case numberOfData > 50000:
		return 12
	case numberOfData > 25000:
		return 13
	case numberOfData > 10000:
		return 14
	case numberOfData > 5000:
		return 15
	case numberOfData > 1000:
		return 16
	default:
		return 17
	}
}

func CheckPermission(refUserOwnerID int, accessCtrl, per1, per2 []string) (bool, bool) {
	isAllData := false
	isOwnerData := false
	if HasPermission(per1, accessCtrl) {
		isAllData = true
	} else if HasPermission(per2, accessCtrl) {
		isOwnerData = true
	} else {
		isAllData = false
		isOwnerData = false
	}

	if isAllData {
		isOwnerData = false
	}
	if refUserOwnerID != 3 {
		isAllData = true
	}
	return isAllData, isOwnerData
}

func CheckPermissionAnalyses(refUserOwnerID int, accessCtrl, per1, per2, per3 []string) (bool, bool) {
	isAllData := false
	isOwnerData := false
	if HasPermission(per1, accessCtrl) {
		isAllData = true
	} else if HasPermission(per2, accessCtrl) {
		isOwnerData = true
	} else {
		isAllData = false
		isOwnerData = false
	}

	if isAllData {
		isOwnerData = false
	}
	if refUserOwnerID != 3 {
		isAllData = true
	}
	return isAllData, isOwnerData
}

func GenerateUUIDV4() string {
	uuid, _ := uuid.NewRandom()
	return uuid.String()
}

func FormatFloatSetPoint(input float64, point int) string {
	setPoint := "%." + fmt.Sprintf(`%d`, point) + "f"
	return fmt.Sprintf(setPoint, input)
}

func SlicesAreEqual(slice1 []int, slice2 []int) bool {
	// Check if the lengths of the slices are the same
	if len(slice1) != len(slice2) {
		return false
	}
	sort.Slice(slice1, func(i, j int) bool {
		return slice1[i] < slice1[j]
	})

	sort.Slice(slice2, func(i, j int) bool {
		return slice2[i] < slice2[j]
	})

	// Check if each element in the slices is the same
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}

	return true
}

func CalculateMapCenter(allmap []string) string {
	var totalLon, totalLat float64
	var pointCount int

	for _, geom := range allmap {
		// Remove "LINESTRING(" and ")" and split by ","
		coords := strings.TrimPrefix(geom, "LINESTRING(")
		coords = strings.TrimSuffix(coords, ")")
		points := strings.Split(coords, ",")

		// Step 3: Loop through each point in the LINESTRING
		for _, point := range points {
			point = strings.TrimSpace(point)
			latLon := strings.Split(point, " ")

			// Parse longitude and latitude
			if len(latLon) == 2 {
				lon, err1 := strconv.ParseFloat(latLon[0], 64)
				lat, err2 := strconv.ParseFloat(latLon[1], 64)

				if err1 == nil && err2 == nil {
					// Accumulate total longitude and latitude
					totalLon += lon
					totalLat += lat
					pointCount++
				}
			}
		}
	}

	// Step 4: Calculate the average (midpoint)
	if pointCount > 0 {
		avgLon := totalLon / float64(pointCount)
		avgLat := totalLat / float64(pointCount)

		// Return the midpoint in LINESTRING format
		return fmt.Sprintf("LINESTRING(%f %f)", avgLon, avgLat)
	}

	// Return an empty string if no points were processed
	return ""
}
