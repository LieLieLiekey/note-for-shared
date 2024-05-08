package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/generative-ai-go/genai"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

const systemPromt = "以该图片为基础，对其进行续集将一个故事或电影，故事背景要符合图片风格（例如有古代神话、魔法、冒险、爱情等），要有故事开头、故事发展、故事高潮和故事结尾"

const geminiAPIKey = "***"

type Response struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func uploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// 允许跨域请求
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 解析 form data
	err := r.ParseMultipartForm(10 << 20) // 限制上传的图片大小为10MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 根据表单字段获取文件
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 读取文件内容
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// 在这里，fileBytes 包含了图片的内容
	// 你可以添加逻辑来处理或分析图片内容
	// 例如，返回图片的简单描述
	// 这里我们假设返回了一个固定的字符串表示已接收到图片
	log.Printf("收到图片，大小：" + fmt.Sprintf("%d", len(fileBytes)) + " 字节\n")

	log.Printf("begin to callGeminiAPI ...\n")
	resp, err := callGeminiAPI(systemPromt, fileBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error, %v", err), http.StatusInternalServerError)
		log.Printf("[ERROR] failed to callGeminiAPI, err: %v\n", err)

		return
	}
	log.Printf("finshed to callGeminiAPI ...\n")

	response := Response{Code: 0, Data: resp.Resp, Msg: "success"}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/upload/image", uploadImageHandler).Methods("POST")
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

type geMiniResp struct {
	Resp string
}

func getTextResultFromResp(v *genai.GenerateContentResponse) string {
	lastResultText := ""

	log.Printf("[DEBUG] getTextResultFromResp begin....\n")

	for _, cd := range v.Candidates {
		log.Printf("[DEBUG] index: %d,\n", cd.Index)

		for _, pat := range cd.Content.Parts {
			switch c := pat.(type) {
			case genai.Text:
				log.Printf("[DEBUG] ----------- role: %s, context text: %s\n", cd.Content.Role, c)

				// 如果是模型的返回，就区最后一个
				if cd.Content.Role == "model" {
					lastResultText = string(c)
				}
			case genai.Blob:
				log.Printf("[DEBUG] ----------- role: %s, context image: c.MIMEType: %s\n", cd.Content.Role, c.MIMEType)
			default:
				continue
			}

		}
	}
	log.Printf("[DEBUG] getTextResultFromResp begin....\n")

	return lastResultText
}

func callGeminiAPI(desc string, imageData []byte) (*geMiniResp, error) {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// For text-and-image input (multimodal), use the gemini-pro-vision model
	model := client.GenerativeModel("gemini-pro-vision")

	prompt := []genai.Part{
		genai.ImageData("jpeg", imageData),
		genai.Text(desc),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return nil, err
	}

	result := getTextResultFromResp(resp)

	bs, _ := json.MarshalIndent(resp, "", "    ")
	fmt.Println(string(bs))

	return &geMiniResp{
		Resp: result,
	}, nil
}
