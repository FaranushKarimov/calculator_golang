package router

import (
	"encoding/json"
	"fmt"
	"hello/internal/models"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func StartRouter() error {
	file, err := os.Open("./config/config.json")
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return err
	}
	var config models.Config

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println(config)

	address := net.JoinHostPort(config.Host, config.Port)
	mux := http.NewServeMux()

	mux.HandleFunc("/calculate", calculate)
	mux.HandleFunc("/history", history)

	log.Fatal(http.ListenAndServe(address, mux))
	return nil
}

var results = []*models.Calculator{}

func calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid method"))
		return
	}

	req := new(models.Calculator)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("json decode err: %+v", err)))
		return
	}

	req.Result = calculator(req.Operation, req.NumberOne, req.NumberSecond)
	results = append(results, req)

	err = json.NewEncoder(w).Encode(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("json encode err: %+v", err)))
		return
	}

	if err := writeToFile(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("write to file err: %+v", err)))
		return
	}
}

func writeToFile(req *models.Calculator) error {
	file, err := os.OpenFile("./internal/history/history.json", os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		return err
	}
	defer file.Close()

	r, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(fmt.Sprintf("%s\r\n", r)))
	if err != nil {
		return err
	}
	return nil
}

func calculator(operation string, firstNum, secondNum float64) float64 {
	switch operation {
	case models.PLUS:
		return firstNum + secondNum
	case models.MINUS:
		return firstNum - secondNum
	case models.MULTIPLY:
		return firstNum / secondNum
	case models.DIVIDE:
		return firstNum * secondNum
	default:
		return 0
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("invalid method"))
	}

	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("json encode err: %+v", err)))
		return
	}
}
