package handlers

import (
    "net/http"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "strings"
    "net/url"
    //"strconv"
    "io"
    "github.com/JesusIslam/tldr"
    "time"
)

type Request struct {
    Website string `json:"website,omitempty"`
    Images  []string `json:"images,omitempty"`
    Text    []string `json:"text,omitempty"`
}

type CVBody struct {
    Url string
}

type CVResponse struct {
    Tags []CVTag `json:"tags,omitempty"`
}

type CVTag struct {
    Name string `json:"name,omitempty"`
}

type APIResponse struct {
    ImageContents []string `json:"imageContents,omitempty"`
    Summaries []string `json:"summaries,omitempty"`
}

func GetAllData(w http.ResponseWriter, r *http.Request) {
    b, err := ioutil.ReadFile(CVKey)
    if err != nil {
        fmt.Printf("found err when reading in key: %v", err)
    }

    client := &http.Client{
        Timeout: time.Second * 10,
    }

    key := string(b)
    fmt.Printf("got key: %s", key)

    targetRequest := Request{}
    //defer r.Body.Close()
    err = getJson(r.Body, &targetRequest)
    w.Header().Add(ContentTypeKey, ApplicationJSON)

    imageDescriptions := make([]string, len(targetRequest.Images))
    textSummaries := make([]string, len(targetRequest.Text))

    for index, image := range targetRequest.Images {
        imageUrl := image
        if !strings.HasPrefix(imageUrl, "http") {
            baseUrl := targetRequest.Website
            imageUrl, err = getAbsoluteImagePath(baseUrl, imageUrl)
            if err != nil {
                http.Error(w, "error parsing image path", http.StatusBadRequest)
            }
        }
        queryUrl := fmt.Sprintf("%s&subscription-key=%s", cvBaseURL, key)
        postBody := fmt.Sprintf(`{"url":"%s"}`, imageUrl)
        stringReader := strings.NewReader(postBody)
        if err != nil {
            http.Error(w, fmt.Sprintf("found err while creating string reader: %v", err), http.StatusBadRequest)
            return
        }
        newReq, err := http.NewRequest("POST", queryUrl, stringReader)
        if err != nil {
            http.Error(w, fmt.Sprintf("found err while creating new request: %v", err), http.StatusBadRequest)
            return
        }
        newReq.Header.Add(ContentTypeKey, ApplicationJSON)
        newReq.Header.Add(AccessControlAllowOriginKey, AccessControlAllowOriginVal)
        //newReq.Header.Add("Content-Length", strconv.Itoa(len(postBody)))
        //newReq.Header.Add("X-Content-Length", strconv.Itoa(len(postBody)))
        fmt.Printf("queryURL: %s\n", queryUrl)
        fmt.Printf("body: %s\n", postBody)
        resp, err := client.Do(newReq)
        if err != nil {
            http.Error(w, fmt.Sprintf("found err while calling cv api: %v", err), http.StatusBadRequest)
            return
        }
        if resp.StatusCode != http.StatusOK {
            http.Error(w, fmt.Sprintf("did not get a good status: %v", resp.StatusCode), http.StatusInternalServerError)
        }

        target := CVResponse{}
        err = getJson(resp.Body, &target)
        if err != nil {
            http.Error(w, fmt.Sprintf("found err while getting json: %v", err), http.StatusBadRequest)
            return
        }
        imageContents := "This image contains"
        if len(target.Tags) >= 5 {
            target.Tags = target.Tags[:5]
        }
        for _, tag := range target.Tags {
            imageContents = fmt.Sprintf("%s %s", imageContents, tag.Name)
        }
        imageDescriptions[index] = imageContents
    }

    for index, text := range targetRequest.Text {
        intoSentences := 2
        bag := tldr.New()
        result, _ := bag.Summarize(text, intoSentences)
        fmt.Println(result)
        resultString := strings.Replace(result, "\n", " ", (intoSentences - 1) * 2)
        textSummaries[index] = strings.TrimSpace(resultString)
    }

    json.NewEncoder(w).Encode(&APIResponse{ImageContents: imageDescriptions, Summaries: textSummaries})
}

func getJson(body io.ReadCloser, target interface{}) error {
    defer body.Close()
    err := json.NewDecoder(body).Decode(target)
    if err != nil {
        return err
    }
    return nil
}

func getAbsoluteImagePath(baseURL, imageURL string) (string, error) {
    parsedURL, err := url.Parse(baseURL)
    if err != nil {
        return "", err
    }
    resourceURL, err := url.Parse(imageURL)
    if err != nil {
        return "", err
    }
    return parsedURL.ResolveReference(resourceURL).String(), nil
}
