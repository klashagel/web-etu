package rest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"

	"strconv"

	"ge.com/fatservice/entity"
	"ge.com/fatservice/interfaces"
	"ge.com/fatservice/websocket"
	"github.com/gorilla/mux"
	"github.com/laurent22/toml-go"
	"github.com/rs/cors"
)

type RestService interface {
	HandleRequests()
}

type restservice struct{}

var (
	//jobMap        *entity.JobMap
	fatController interfaces.IFatController
	settings      toml.Document
)

func NewRestService(ctrl interfaces.IFatController, tdoc toml.Document) RestService {
	fatController = ctrl
	settings = tdoc
	return &restservice{}
}

func (*restservice) HandleRequests() {
	log.Println("Rest service initialized")
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/etu/getallcontrollers", returnAllControllers)
	myRouter.HandleFunc("/etu/getinfotext/{ip}", getControllerInfoText)
	myRouter.HandleFunc("/etu/getcontrollerdata/{ip}", getControllerData)
	myRouter.HandleFunc("/etu/getdiscreteinputs/{ip}", getDiscreteInputs)
	myRouter.HandleFunc("/etu/getcoils/{ip}", getCoils)
	myRouter.HandleFunc("/api/epic4updatefirmware", epic4UpdateFirmware).Methods("PUT")
	myRouter.HandleFunc("/api/epic4burnimage", epic4BurnImage).Methods("PUT")
	myRouter.HandleFunc("/api/epic4setmac", epic4SetMacAddress).Methods("PUT")

	// Test cases
	myRouter.HandleFunc("/testcase/create", createTestCase).Methods("POST")
	myRouter.HandleFunc("/testcase/jobs/{id}", getTestCaseJobs)
	myRouter.HandleFunc("/testcase/delete/{id}", deleteTestCase).Methods("DELETE")
	myRouter.HandleFunc("/testcase/update", updateTestCase).Methods("PUT")

	// Jobs
	myRouter.HandleFunc("/job/createreport/{jbid}/{rptid}", createReport)
	myRouter.HandleFunc("/job/getrunning", getAllJobs)
	myRouter.HandleFunc("/job/getmap", getJobMap)
	myRouter.HandleFunc("/job/cancelall", cancelAllJobs)
	myRouter.HandleFunc("/job/cancel", cancelJob).Methods("PUT")
	myRouter.HandleFunc("/job/nextstate", nextStateJob).Methods("PUT")
	//myRouter.HandleFunc("/job/search/{sc}", filterJobs)
	myRouter.HandleFunc("/job/status/{id}", getJobStatus)
	myRouter.HandleFunc("/job/start/{id}", startJob)
	myRouter.HandleFunc("/job/pause/{id}", pauseJob).Methods("PUT")
	myRouter.HandleFunc("/job/getinfotext/{id}", getJobInfoText)
	myRouter.HandleFunc("/job/getstatustext/{id}", getJobStatusText)
	myRouter.HandleFunc("/job/getjoblog/{id}", getJobLog)
	myRouter.HandleFunc("/job/archive/{id}", archiveJob)
	myRouter.HandleFunc("/job/create", createJob).Methods("POST")
	myRouter.HandleFunc("/job/update", updateJob).Methods("PUT")
	myRouter.HandleFunc("/job/delete/{id}", deleteJob).Methods("DELETE")
	myRouter.HandleFunc("/job/start", startJob).Methods("PUT")
	myRouter.HandleFunc("/jobtype/get", getJobTypes)
	myRouter.HandleFunc("/job/resetjob", resetJob).Methods("PUT")
	myRouter.HandleFunc("/job/restart", restartJob).Methods("PUT")
	// Articles
	myRouter.HandleFunc("/articles/get", getArticles)
	myRouter.HandleFunc("/article/testcase/{id}", getTestCases)
	myRouter.HandleFunc("/articles/getpage/{page}/{rec}", getArticlePage)
	myRouter.HandleFunc("/article/create", createArticle).Methods("POST")
	myRouter.HandleFunc("/article/update", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/delete/{id}", deleteArticle).Methods("DELETE")

	// Helper registers

	myRouter.HandleFunc("/outputtype/get", getOutputTypes)
	myRouter.HandleFunc("/sirtype/get", getSirTypes)
	myRouter.HandleFunc("/reports/get", getFatReports)

	// Keypad data
	myRouter.HandleFunc("/serial/digital/getall", getSerialDigitalData)
	myRouter.HandleFunc("/serial/digital/write", writeSerial).Methods("PUT")
	myRouter.HandleFunc("/api/serial/write", writeSerialCommand).Methods("PUT")
	myRouter.HandleFunc("/serial/analog/getall", getSerialAnalogData)

	myRouter.HandleFunc("/modbus/write", writeModbus).Methods("PUT")

	myRouter.HandleFunc("/lua/run/{fileId}", runLua).Methods("PUT")
	myRouter.HandleFunc("/lua/cancel", cancelLua).Methods("PUT")
	myRouter.HandleFunc("/lua/results", getLuaResults).Methods("GET")
	myRouter.HandleFunc("/lua/files", getLuaFilesInFolder).Methods("GET")

	myRouter.HandleFunc("/api/files/{fileId}", FetchLuaCodeHandler).Methods("GET")
	myRouter.HandleFunc("/api/files/{fileId}", SaveLuaCodeHandler).Methods("POST")
	myRouter.HandleFunc("/api/check-lua-syntax", CheckLuaSyntaxHandler).Methods("POST")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", //or you can your header key values which you are using in your application

		},
	})

	handler := corsOpts.Handler(myRouter)
	log.Fatal(http.ListenAndServe(settings.GetString("rest.port"), handler))

	//log.Fatal(http.ListenAndServe(":12345", handlers.CORS(originsOk, headersOk, methodsOk)(myRouter)))
}

func getNormalizedFilePath(configPath, fileName string) string {
	// Normalize the path to use slashes
	normalizedPath := filepath.FromSlash(configPath)

	// Join the normalized path with the file name
	return filepath.Join(normalizedPath, fmt.Sprintf("%s", fileName))
}

func getLuaFilesInFolder(w http.ResponseWriter, r *http.Request) {
	folderPath := settings.GetString("etuserver.luacodepath") // Specify the folder path

	files, err := os.ReadDir(folderPath)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	var filenames []string
	for _, file := range files {
		if !file.IsDir() {
			filenames = append(filenames, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filenames)
}

// FetchLuaCodeHandler handles GET requests to fetch Lua code based on fileId
func FetchLuaCodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["fileId"]

	// For demo purposes, assume Lua code is stored in a file named <fileId>.lua
	filePath := getNormalizedFilePath(settings.GetString("etuserver.luacodepath"), fileName)
	fmt.Printf("FetchLuaCodeHandler %s\n", filePath)

	// Read the file
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Respond with the code
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"code": string(code)})
}

// SaveLuaCodeHandler handles POST requests to save Lua code
func SaveLuaCodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SaveLuaCodeHandler called")
	vars := mux.Vars(r)
	fileName := vars["fileId"]

	fmt.Println("File id ", fileName)

	// Parse the request body
	var requestBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For demo purposes, save Lua code to a file named <fileId>.lua
	filePath := getNormalizedFilePath(settings.GetString("etuserver.luacodepath"), fileName)
	fmt.Printf("SaveLuaCodeHandler %s\n", filePath)

	// Write the file
	err := ioutil.WriteFile(filePath, []byte(requestBody.Code), 0644)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CheckLuaSyntaxHandler handles POST requests to check Lua code syntax
func CheckLuaSyntaxHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestBody struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a temporary file to store the Lua code
	tempFile, err := ioutil.TempFile("", "*.lua")
	if err != nil {
		http.Error(w, "Failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	// Write the Lua code to the temporary file
	if _, err := tempFile.WriteString(requestBody.Code); err != nil {
		http.Error(w, "Failed to write to temp file", http.StatusInternalServerError)
		return
	}
	tempFile.Close()

	// Execute the Lua interpreter to check syntax
	cmd := exec.Command("lua", "-p", tempFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If syntax error, output will contain the error message
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": string(output)})
		return
	}

	// If no errors, report syntax is correct
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Syntax is correct"})
}

func resetJob(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Update job called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	go fatController.ResetJob(job)

	fatController.CancelJob(job.Id)

}

func restartJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Restartjob called")
	reqBody, _ := io.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	fatController.ResetJob(job)
	fatController.CancelJob(job.Id)
	fatController.StartJob(job.Id)
}

func archiveJob(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rest: Archive job called")
	vars := mux.Vars(r)
	key := vars["id"]
	id, _ := strconv.Atoi(key)
	job := fatController.GetJob(id)
	//fmt.Println("Rest: Got job with startdate and enddate", job.StartTime, job.EndTime)
	fatController.ArchiveJob(job)
}

func getJobLog(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Rest: Get job log called")
	vars := mux.Vars(r)
	key := vars["id"]
	id, _ := strconv.Atoi(key)
	//job := fatController.GetJob(id)

	log := fatController.GetJobLog(id)
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if log != nil {
		json.NewEncoder(w).Encode(log)
	}

}

func createReport(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Create report called")

	vars := mux.Vars(r)
	jbId, _ := strconv.Atoi(vars["jbid"])
	rptId, _ := strconv.Atoi(vars["rptid"])

	job := fatController.GetJob(jbId)
	//fmt.Println(job.Id)
	//outfile := fmt.Sprintf("hvreport_%v.pdf", job.Id)

	reportscript := settings.GetString("report.reportscript")
	//outdir := "/home/ge/sambashare/"

	// rendercmd := fmt.Sprintf("rmarkdown::render('%s', params = list(jobid = '%d', reportid = '%d'),output_file='%s', output_dir ='%s')", reportdef, jbId, rptId, outfile, outdir)

	//rendercmd := fmt.Sprintf("%s  %d %d", reportdef, jbId, rptId)

	cmd := exec.Command("Rscript", reportscript, strconv.Itoa(jbId), strconv.Itoa(rptId))
	fmt.Println(cmd)

	//Rscript-e rmarkdown::render('/home/ge/sambashare/first.Rmd', params = list(jobid = '73'),output_file='hvreport_73.pdf', output_dir ='/home/ge/sambashare/')

	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	// go func() {
	//     err := <-ch

	//     fmt.Println("returned err", err)
	//     websocket.BroadCast(entity.BroadCastMessage{
	//         EventId: entity.Job_ReportFinished,
	//         Data:    strconv.Itoa(job.Id),
	//     })

	//     return
	// }()

	go func() {
		err := <-ch

		fmt.Println("returned err", err)
		websocket.BroadCast(entity.BroadCastMessage{
			EventId: entity.Job_ReportFinished,
			Data:    strconv.Itoa(job.Id),
		})

	}()
}

func returnAllControllers(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint Hit: returnAllArticles")
	// Needed for flutter restapi cors localhost...
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	json.NewEncoder(w).Encode(fatController.GetActiveControllers())

}

func getJobInfoText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	jbid, err := strconv.Atoi(key)

	if err != nil {
		http.Error(w, "Controller not responding", 500)
		return
	}
	//fmt.Println("getJobInfoText called for ", jbid)
	infoText := fatController.GetJobInfoText(jbid)

	json.NewEncoder(w).Encode(&entity.StringResp{Data: infoText})
}

func getJobStatusText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	jbid, err := strconv.Atoi(key)

	if err != nil {
		http.Error(w, "Controller not responding", 500)
		return
	}
	//fmt.Println("getJobStatusText called for ", jbid)
	infoText := fatController.GetJobStatusText(jbid)

	json.NewEncoder(w).Encode(&entity.StringResp{Data: infoText})
}

func getControllerInfoText(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ip := vars["ip"]
	//fmt.Println("getControllerInfoTextm called for ", ip)
	infoText := fatController.GetControllerInfoText(ip)

	json.NewEncoder(w).Encode(&entity.StringResp{Data: infoText})
}

func getControllerData(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ip := vars["ip"]
	//fmt.Println("getControllerInfoTextm called for ", ip)
	data := fatController.GetControllerData(ip)

	json.NewEncoder(w).Encode(data)
}

func getDiscreteInputs(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ip := vars["ip"]
	//fmt.Println("getControllerInfoTextm called for ", ip)
	data := fatController.GetControllerDiscreteInputs(ip)

	json.NewEncoder(w).Encode(data)
}

func getCoils(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ip := vars["ip"]
	//fmt.Println("getControllerInfoTextm called for ", ip)
	data := fatController.GetControllerCoils(ip)

	json.NewEncoder(w).Encode(data)
}

func getJobTypes(w http.ResponseWriter, r *http.Request) {

	jobtypes := fatController.GetJobTypes()
	json.NewEncoder(w).Encode(jobtypes)
}

func getArticles(w http.ResponseWriter, r *http.Request) {

	articles := fatController.GetArticles()
	json.NewEncoder(w).Encode(articles)
}

func getOutputTypes(w http.ResponseWriter, r *http.Request) {

	outputType := fatController.GetOutputTypes()
	json.NewEncoder(w).Encode(outputType)
}

func getSirTypes(w http.ResponseWriter, r *http.Request) {

	sirTypes := fatController.GetSirTypes()
	json.NewEncoder(w).Encode(sirTypes)
}

func getFatReports(w http.ResponseWriter, r *http.Request) {

	output := fatController.GetFatReports()
	json.NewEncoder(w).Encode(output)
}

func startJob(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Start job called")
	vars := mux.Vars(r)
	jbIdStr := vars["id"]

	jbId, _ := strconv.Atoi(jbIdStr)

	fatController.StartJob(jbId)

}

func pauseJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	jbIdStr := vars["id"]

	jbId, _ := strconv.Atoi(jbIdStr)

	fmt.Println("Pause job called ", jbId)
	fatController.PausJob(jbId)

}

/* func stopJob(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Stop job called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	fatController.RunJob(job)

	//fmt.Fprintf(w, "Article created")

} */

func getJobStatus(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]
	id, err := strconv.Atoi(key)

	if err != nil {
		fmt.Println("error in job key")
	}

	if val, ok := fatController.GetJobMap().Jobs[id]; ok {
		val.Stats.M.Lock()
		ret := &entity.JobStats{
			State:         val.Stats.State,
			Status:        val.Stats.Status,
			StateName:     val.StateDef.Exec[val.Stats.State],
			DataPoints:    val.Stats.DataPoints,
			RemainingTime: val.Stats.RemainingTime,
			Tags:          val.Stats.Tags,
			Fields:        val.Stats.Fields,
			ReadErrors:    val.Stats.ReadErrors,
			Failed:        val.Stats.Failed,
		}

		json.NewEncoder(w).Encode(ret)
		val.Stats.M.Unlock()
	}

}

func getArticlePage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	page := vars["page"]
	rec := vars["rec"]

	iPage, _ := strconv.Atoi(page)
	iRec, _ := strconv.Atoi(rec)

	articles, _ := fatController.GetArticlePage(iPage, iRec)
	json.NewEncoder(w).Encode(articles)

}

func createTestCase(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ts entity.TestCase
	json.Unmarshal(reqBody, &ts)

	uid := r.Header.Get("uid")

	fmt.Println("createtestcase: Got req header name :", uid)

	id, err := fatController.CreateTestCase(uid, ts)
	if err != nil {
		log.Println("Error creating testcase", err)
	}
	//msg := entity.BroadCastMessage{EventId: entity.Article_Create, Data: article.SerialNumber}
	//go websocket.BroadCast(msg)

	fmt.Fprintf(w, "%v", id)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article entity.Article
	json.Unmarshal(reqBody, &article)

	uid := r.Header.Get("uid")

	fmt.Println("createArticle: Got req header name :", uid)

	err := fatController.CreateArticle(uid, article)
	if err != nil {
		log.Println("Error creating article", err)
	}
	//msg := entity.BroadCastMessage{EventId: entity.Article_Create, Data: article.SerialNumber}
	//go websocket.BroadCast(msg)

	fmt.Fprintf(w, "Article created")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	id, _ := strconv.Atoi(key)

	uid := r.Header.Get("uid")
	//fmt.Println("Got req header name :", ua)

	err := fatController.DeleteArticle(uid, id)

	/* 	if err == nil {
		msg := entity.BroadCastMessage{EventId: entity.Article_Delete, Data: ua}
		go websocket.BroadCast(msg)
	} */
	_ = err
	fmt.Fprintf(w, "Article deleted")

}

func deleteJob(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	id, _ := strconv.Atoi(key)
	err := fatController.DeleteJob(id)

	_ = err
	fmt.Fprintf(w, "Article deleted")

}

func deleteTestCase(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["id"]

	id, _ := strconv.Atoi(key)
	err := fatController.DeleteTestCase(id)

	_ = err
	if err != nil {
		http.Error(w, "Error deleteing article", 400)
	} else {
		//http.Error(w, "Error deleteing article", 400)
		w.WriteHeader(200)
	}

}

func createJob(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	id, err := fatController.CreateJob(job)
	if err != nil {
		log.Println("Error creating job", err)
	}

	//msg := entity.BroadCastMessage{EventId: entity.Job_Create, Data: "No message"}
	//go websocket.BroadCast(msg)

	fmt.Fprintf(w, "%v", id)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update article called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var art entity.Article
	json.Unmarshal(reqBody, &art)

	fatController.UpdateArticle(art)
	//fmt.Println(job)
}

func updateTestCase(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update testcase called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ts entity.TestCase
	json.Unmarshal(reqBody, &ts)
	fatController.UpdateTestCase(ts)
	fmt.Println(ts)
}

func updateJob(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Update job called")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	fatController.UpdateJob(job)
	//fmt.Println(job)
}

func cancelAllJobs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cancel all jobs called")
	for k := range fatController.GetJobMap().Jobs {
		fatController.GetJobMap().Jobs[k].Cancel()
	}
	fmt.Fprintf(w, "Closing all jobs")
}

func nextStateJob(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Executing next state for job")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)
	fmt.Println("Executing next state for job", job.Id)
	go fatController.ExecNextStateJob(job.Id)

	fmt.Fprintf(w, "Exec next state")
}

func cancelJob(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var job entity.Job
	json.Unmarshal(reqBody, &job)

	fatController.CancelJob(job.Id)

	fmt.Fprintf(w, "Close job")
}

func getJobMap(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	ret := fatController.GetJobMap()

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(ret)

}

func getAllJobs(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	ret := make([]*entity.Job, 0)
	for _, v := range fatController.GetJobMap().Jobs {

		ret = append(ret, v.Job)
	}

	// Need to sort otherwize reciving list/table will jumble on update
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Id < ret[j].Id
	})

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	json.NewEncoder(w).Encode(ret)

}

func getTestCaseJobs(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tcid := vars["id"]
	//fmt.Println("Rest: Filter job called ", sc)
	jobs := fatController.GetTestCaseJobs(tcid)

	json.NewEncoder(w).Encode(jobs)

}

func getTestCases(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	//fmt.Println("Rest: Filter job called ", sc)
	ts := fatController.GetTestCases(id)

	json.NewEncoder(w).Encode(ts)

}

func getSerialDigitalData(w http.ResponseWriter, r *http.Request) {
	digitalData := fatController.GetSerialDigitalData()
	//fmt.Println(digitalData)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(digitalData)
}

func getSerialAnalogData(w http.ResponseWriter, r *http.Request) {
	analogData := fatController.GetSerialAnalogData()
	//fmt.Println(digitalData)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analogData)
}

type ButtonState struct {
	ButtonNumber int  `json:"buttonNumber"`
	NewState     bool `json:"newState"`
}

type SerialAction struct {
	Command string `json:"command"`
}

type ModbusWriteRec struct {
	Ip       string `json:"ip"`
	Register string `json:"register"`
	Value    int    `json:"value"`
}

func writeSerial(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Writeserial called")
	reqBody, _ := io.ReadAll(r.Body)
	//fmt.Println(string(reqBody))

	var buttonState ButtonState

	// Unmarshal the request body into the struct
	if err := json.Unmarshal(reqBody, &buttonState); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Print the button state
	//fmt.Printf("Button %d has been toggled to %t\n", buttonState.ButtonNumber, buttonState.NewState)

	fatController.SetDigitalIO(buttonState.ButtonNumber, buttonState.NewState)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "PUT request body received successfully")
}

func writeSerialCommand(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	var rec SerialAction
	// Unmarshal the request body into the struct
	if err := json.Unmarshal(reqBody, &rec); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Print the button state
	fmt.Printf("Write analog %s\n", rec.Command)

	fatController.WriteAnalog(rec.Command)

	response := map[string]string{"message": "PUT request body received successfully"}

	// Set the response headers and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func writeModbus(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Writemodbus called")
	reqBody, _ := io.ReadAll(r.Body)
	fmt.Println(string(reqBody))

	var rec ModbusWriteRec

	// Unmarshal the request body into the struct
	if err := json.Unmarshal(reqBody, &rec); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Print the button state
	fmt.Printf("Modbus write to %s,  in register %s value %d\n", rec.Ip, rec.Register, rec.Value)

	fatController.ModbusWrite(rec.Ip, rec.Register, uint16(rec.Value))

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "PUT request body received successfully")
}

func runLua(w http.ResponseWriter, r *http.Request) {
	// Extract the fileId parameter from the URL
	vars := mux.Vars(r)
	fileIdStr, ok := vars["fileId"]
	if !ok {
		http.Error(w, "Missing fileId parameter", http.StatusBadRequest)
		return
	}

	// Log the extracted fileIdStr for debugging
	fmt.Println("Run lua called with fileId:", fileIdStr)

	// Read the JSON body from the request
	var params map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Log the received parameters for debugging
	fmt.Printf("Received parameters: %v\n", params)

	// Extract the ipaddress from the params map (if not already present, set it)
	ipaddress, ipExists := params["ipaddress"].(string)
	if !ipExists || ipaddress == "" {
		http.Error(w, "Missing or invalid 'ipaddress' parameter", http.StatusBadRequest)
		return
	}

	// Ensure ipaddress is in the params map
	params["ipaddress"] = ipaddress

	// Log the extracted IP address for debugging
	fmt.Println("Received IP address:", ipaddress)

	// Process the fileIdStr and params as needed
	filePath := getNormalizedFilePath(settings.GetString("etuserver.luacodepath"), fileIdStr)

	// Pass the file path and parameters (including IP) to RunLuaScript
	err := fatController.RunLuaScript(filePath, params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error running Lua script: %v", err), http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

func cancelLua(w http.ResponseWriter, r *http.Request) {

	fatController.CancelLuaScript()
	w.WriteHeader(http.StatusOK)
}

func getLuaResults(w http.ResponseWriter, r *http.Request) {
	results := fatController.GetLuaResults()
	resultsString := fmt.Sprintf("%v", results)
	fmt.Printf("GetLuaResults called %s", resultsString)

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Attempt to encode the results to JSON
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		// If there is an error encoding the JSON, respond with an error status
		http.Error(w, "Failed to encode results as JSON", http.StatusInternalServerError)
		return
	}

	// No need to call WriteHeader explicitly, as Encode sets the status to 200 OK by default
}

func epic4UpdateFirmware(w http.ResponseWriter, r *http.Request) {
	// Set headers for streaming the response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Retrieve the IP address from the headers
	ip := r.Header.Get("ip")
	if ip == "" {
		http.Error(w, "IP header missing", http.StatusBadRequest)
		return
	}

	// Determine which batch file to execute and its arguments
	batFile := settings.GetString("iwave.updatefirmware_batfile")
	var args []string

	// Example logic to determine if arguments are needed

	args = []string{ip} // Arguments for install.bat

	// Execute the batch file with arguments
	executeBatchFile(w, batFile, args...)
}

func epic4BurnImage(w http.ResponseWriter, r *http.Request) {

	// Set the header for streaming the response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	batFile := settings.GetString("iwave.burnimage_batfile")

	fmt.Println("Batfile ", batFile)

	executeBatchFile(w, batFile)
}

type MacAddressRequest struct {
	MacAddress1 string `json:"macAddress1"`
	MacAddress2 string `json:"macAddress2"`
	ComPort     string `json:"comPort"`
}

func epic4SetMacAddress(w http.ResponseWriter, r *http.Request) {
	// Set headers for streaming the response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Parse the JSON request body
	var requestBody MacAddressRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract MAC addresses
	macAddress1 := requestBody.MacAddress1
	macAddress2 := requestBody.MacAddress2
	comPort := requestBody.ComPort

	// Get the path to the batch file
	batFile := settings.GetString("iwave.setmac_batfile")

	// Log extracted MAC addresses and batch file path
	fmt.Printf("MAC Address 1: %s\n", macAddress1)
	fmt.Printf("MAC Address 2: %s\n", macAddress2)
	fmt.Printf("Batch file: %s\n", batFile)

	// Execute the batch file with MAC addresses as arguments
	args := []string{comPort, macAddress1, macAddress2}
	executeBatchFile(w, batFile, args...)
}
func executeBatchFile(w http.ResponseWriter, batFile string, args ...string) {
	fmt.Println("Batch file path:", batFile)
	fmt.Println("Arguments:", args)

	// Extract the directory from the batch file path
	workingDir := filepath.Dir(batFile)
	fmt.Println("Extracted working directory:", workingDir)

	// Create the command with the batch file and arguments
	cmd := exec.Command("cmd.exe", append([]string{"/C", batFile}, args...)...)
	cmd.Dir = workingDir // Set the working directory for the command

	// Create pipes to capture the output and error streams
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, "Error creating stdout pipe", http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, "Error creating stderr pipe", http.StatusInternalServerError)
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		http.Error(w, "Error starting command", http.StatusInternalServerError)
		return
	}

	// Combine stdout and stderr
	combined := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(combined)

	// Stream the output to the client
	for scanner.Scan() {

		websocket.BroadCast(entity.BroadCastMessage{
			EventId: entity.Cmd_Message,
			Data:    scanner.Text(),
		})
	}

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		websocket.BroadCast(entity.BroadCastMessage{
			EventId: entity.Cmd_Message,
			Data:    fmt.Sprintf("Command finished with error: %v\n\n", err),
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Batch file execution failed.\n"))
		return
	}

	// Check the exit code
	exitCode := cmd.ProcessState.ExitCode()
	if exitCode == 0 {
		websocket.BroadCast(entity.BroadCastMessage{
			EventId: entity.Cmd_Message,
			Data:    "Command finished successfully.\n\n",
		})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Batch file executed successfully.\n"))
	} else {
		websocket.BroadCast(entity.BroadCastMessage{
			EventId: entity.Cmd_Message,
			Data:    "Command finished with errors.\n\n",
		})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Batch file execution resulted in errors.\n"))
	}
}
