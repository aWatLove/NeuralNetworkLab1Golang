package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type NeuralNetwork struct {
	Neurons [][]float64   `json:"neurons"` // оутпуты на нейроне и сами нейроны
	W       [][][]float64 `json:"w"`       // веса [слой][нейрон][вес]
	LR      float64       `json:"lr"`      // скорость обучения
	EPOCH   int           `json:"epoch"`   // кол-во эпох обучения
	Dw      [][][]float64 `json:"dw"`      // дельты весов // неиспользуется // -deprecated
	Mu      float64       `json:"mu"`      // коэф. инерционности
	T       int           `json:"t"`       // номер - текущей итерации
}

func NewNeuralNetwork(neurons [][]float64, w [][][]float64, LR float64, EPOCH int, dw [][][]float64, mu float64) *NeuralNetwork {
	return &NeuralNetwork{Neurons: neurons, W: w, LR: LR, EPOCH: EPOCH, Dw: dw, Mu: mu, T: 1}
}

var dicts = make(map[string]float64) // словарь
var dictsInd float64                 // index словаря

func dictAdd(s string) (i float64) {
	if e, ok := dicts[s]; !ok {
		dictsInd++
		dicts[s] = dictsInd
		i = dictsInd
	} else {
		i = e
	}
	return i
}

var results = [3]string{"Надежный клиент", "Клиент с высоким риском", "Ненадежный клиент"}

const a = 1 // a - альфа для функции активации
var dict = make(map[string]int)

func main() {
	trainData, expRes, err := loadDataSet()
	if err != nil {
		log.Fatalln(err)
	}

	inputData := len(trainData[0])

	neurons := createNeurons(inputData, 12, 10, 8, 6, 3)
	w := createWeights(neurons)
	generateWeights(w)
	dw := createWeights(neurons)
	var nn = NeuralNetwork{Neurons: neurons, W: w, LR: 0.2, EPOCH: 100000, Dw: dw, Mu: 0.1}

	train(&nn, trainData, expRes)

	nn.saveNN("nn")

	predict(&nn, trainData[0])
	out := imvia(nn.Neurons[len(nn.Neurons)-1])
	fmt.Println(nn.Neurons[len(nn.Neurons)-1])
	fmt.Println("result:", results[out])

	//fmt.Printf("NN:\n %v\n", nn.neurons)

	//predict(&nn, trainData[61])
	//fmt.Printf("NN:\n %v\n", nn.neurons)
}

func createNeurons(a ...int) (neurons [][]float64) {
	neurons = make([][]float64, len(a))
	for i, c := range a {
		neurons[i] = make([]float64, c)
	}
	return neurons
}

func createMatrixByNN(nn [][]float64) (m [][]float64) {
	m = make([][]float64, len(nn))
	for i, l := range nn {
		m[i] = make([]float64, len(l))
	}
	return m
}

func createWeights(neurons [][]float64) (w [][][]float64) {
	w = make([][][]float64, len(neurons)-1)
	for i := 0; i < len(neurons)-1; i++ {
		w[i] = make([][]float64, len(neurons[i]))
		for j := 0; j < len(neurons[i]); j++ {
			w[i][j] = make([]float64, len(neurons[i+1]))
		}
	}
	return w
}

func minMaxVal(data [][]float64, min []float64, max []float64) {
	for i := 3; i < len(data[0]); i++ {
		min[i-3] = data[0][i]
		max[i-3] = data[0][i]
	}

	for i := 0; i < len(data); i++ {
		for j := 3; j < len(data[i]); j++ {
			if data[i][j] < min[j-3] {
				min[j-3] = data[i][j]
			}
			if data[i][j] > max[j-3] {
				max[j-3] = data[i][j]
			}
		}
	}
}

func normilizeData(data [][]float64) {
	minV := make([]float64, 7)
	maxV := make([]float64, 7)
	minMaxVal(data, minV, maxV)
	fmt.Println(minV, maxV)
	for i := 0; i < len(data); i++ {
		for j := 3; j < len(data[i]); j++ {
			data[i][j] = (data[i][j] - minV[j-3]) / (maxV[j-3] - minV[j-3])
		}
	}

}

// загрузка data set'а из файла
func loadDataSet() (trainData [][]float64, expResults [][]float64, err error) {
	// Открываем CSV файл
	file, err := os.Open("dataset.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}
	defer file.Close()

	// Создаем CSV Reader
	reader := csv.NewReader(file)

	// Читаем все записи из CSV файла
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	}

	//

	// Проходим по всем записям и преобразуем числовые значения в формат float64
	var data [][]float64
	for _, row := range records[1:] {
		var vector []float64
		ind := dictAdd(row[0])
		switch ind {
		case 0:
			vector = append(vector, 1)
			vector = append(vector, 0)
			vector = append(vector, 0)
		case 1:
			vector = append(vector, 0)
			vector = append(vector, 1)
			vector = append(vector, 0)
		case 2:
			vector = append(vector, 0)
			vector = append(vector, 0)
			vector = append(vector, 1)
		default:
			vector = append(vector, 0)
			vector = append(vector, 0)
			vector = append(vector, 0)
		}
		for _, value := range row[1 : len(row)-1] { // Начинаем с 2-го элемента, так как первые два это строковые значения
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			vector = append(vector, floatValue)
		}
		data = append(data, vector)

		switch row[len(row)-1] {
		case "Надежный клиент":
			expResults = append(expResults, []float64{1, 0, 0})
		case "Клиент с высоким риском":
			expResults = append(expResults, []float64{0, 1, 0})
		case "Ненадежный клиент":
			expResults = append(expResults, []float64{0, 0, 1})
		}
	}
	normilizeData(data)
	trainData = data

	return trainData, expResults, nil
}

func loadTestData() ([][]float64, error) { //TODO: доделать!

	// Открываем CSV файл
	file, err := os.Open("dataset.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	defer file.Close()

	// Создаем CSV Reader
	reader := csv.NewReader(file)

	// Читаем все записи из CSV файла
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	fmt.Printf("%v\n\n\n", records)
	// Проходим по всем записям и преобразуем числовые значения в формат float64
	var data [][]float64
	for _, row := range records[1:] {
		var vector []float64
		ind := dictAdd(row[0])
		vector = append(vector, ind)
		for _, value := range row[1:] { // Начинаем с 2-го элемента, так как первые два это строковые значения
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			vector = append(vector, floatValue)
		}
		data = append(data, vector)
	}

	// Используйте данные как входные данные для вашей нейронной сети
	fmt.Println(data)

	return nil, nil
}

func (nn *NeuralNetwork) saveNN(filename string) error {
	jsonData, err := json.MarshalIndent(nn, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Данные нейронной сети успешно сохранены в файл %s\n", filename)
	return nil
}

func loadNN(filename string) (*NeuralNetwork, error) {
	var nn NeuralNetwork
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileData, &nn)
	if err != nil {
		return nil, err
	}
	return &nn, nil
}

func saveWeights(filename string, w [][][]float64) { // сохранить мозги (веса) в файл
	jsonData, err := json.Marshal(w)
	if err != nil {
		fmt.Println("Ошибка сохранения весов в JSON-файл", err)
	}
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка записи весов в JSON-файл", err)
	}
	fmt.Printf("Данные успешно сохранены в файл %s\n", filename)
}

func loadWeights(filename string) ([][][]float64, error) { // загрузить пресет мозгов (весов) // todo
	var w [][][]float64
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(fileData, &w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func generateWeights(w [][][]float64) {
	for i := 0; i < len(w); i++ {
		for j := 0; j < len(w[i]); j++ {
			for k := 0; k < len(w[i][j]); k++ {
				// Инициализация весов методом "He"
				w[i][j][k] = rand.NormFloat64() * math.Sqrt(2.0/float64(len(w[i][j])))
			}
		}
	}
}

func activate(s float64) float64 {
	return (1 / (1 + math.Pow(math.E, -s)))
}

func forward(nn *NeuralNetwork) { // функция прямого распространения
	for n := 0; n < len(nn.W); n++ { // цикл по слоям нейронов
		for i, _ := range nn.Neurons[n+1] {
			var res float64 = 0
			for j, input := range nn.Neurons[n] {
				res += input * nn.W[n][j][i] // i - к output нейрону, j - из input нейрона
			}
			nn.Neurons[n+1][i] = activate(res)
		}
	}
}

func backProp(nn *NeuralNetwork, exp []float64) {
	forward(nn)
	lastLayer := len(nn.Neurons) - 1
	m := createMatrixByNN(nn.Neurons)

	// Вычисляем градиенты в последнем слое
	for i, n := range nn.Neurons[lastLayer] {
		m[lastLayer][i] = n * (1 - n) * (n - exp[i])
	}

	// Обратное распространение ошибки через скрытые слои
	for i := lastLayer - 1; i > 0; i-- {
		for j := 0; j < len(nn.Neurons[i]); j++ {
			var sum float64
			for k, elem := range m[i+1] {
				sum += elem * nn.W[i][j][k]
			}
			m[i][j] = sum * (1 - nn.Neurons[i][j])
		}
	}

	for i, wi := range nn.W {
		for j, wij := range wi {
			for k, _ := range wij {
				deltaW := -(nn.LR * m[i+1][k] * nn.Neurons[i][j])
				nn.W[i][j][k] += deltaW
				//deltaW := -(nn.LR * (nn.mu*nn.dw[i][j][k]*(float64(nn.t)-1) + (1-nn.mu)*m[i+1][k]*nn.neurons[i][j]))
				//nn.dw[i][j][k] = deltaW
				//nn.t += 1
				//nn.w[i][j][k] += deltaW
			}
		}
	}
}

func train(nn *NeuralNetwork, data [][]float64, exp [][]float64) {
	maxAcc := 0.0
	var niceW = createWeights(nn.Neurons)
	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам
		for d := 0; d < len(data); d++ { // цикл по дата сету
			for i := 0; i < len(data[d]); i++ {
				nn.Neurons[0][i] = data[d][i] // присваиваем входным нейронам данные из дата-сета
			}
			backProp(nn, exp[d]) // вычисляем ошибку и корректируем веса
		}
		accuracy, errExp := evaluate(nn, data, exp)

		fmt.Printf("Epoch: %d, Accuracy: %.6f%%, ResExp: %f\n", e+1, accuracy*100, errExp)

		if accuracy > maxAcc { // копирую максимально точные веса
			maxAcc = accuracy
			for i := 0; i < len(niceW); i++ {
				for j := 0; j < len(niceW[i]); j++ {
					for k := 0; k < len(niceW[i][j]); k++ {
						niceW[i][j][k] = nn.W[i][j][k]
					}
				}
			}
		}

	}
	saveWeights(fmt.Sprintf("w acc %.f", maxAcc*100), niceW) // save best weights
	nn.W = niceW
}

func predict(nn *NeuralNetwork, data []float64) []float64 { // вычислить
	for i, _ := range nn.Neurons[0] {
		nn.Neurons[0][i] = data[i]
	}
	forward(nn)

	ll := len(nn.Neurons) - 1 // last layer
	outputs := make([]float64, len(nn.Neurons[ll]))
	for i, _ := range nn.Neurons[ll] {
		outputs[i] = nn.Neurons[ll][i]
	}
	return outputs
}

func evaluate(nn *NeuralNetwork, data [][]float64, exp [][]float64) (accuracy float64, errSum float64) {
	var count int
	var res float64

	for i, d := range data {
		outputs := predict(nn, d)
		o := imvia(outputs) // output
		t := imvia(exp[i])  // target - ожидаемое значение
		if o == t {
			count++
		}
		for j, e := range outputs {
			res += math.Pow((e - exp[i][j]), 2)
		}
		errSum = res * 0.5

	}
	accuracy = float64(count) / float64(len(data))
	errSum = errSum / float64(len(exp))
	return accuracy, errSum
}

func imvia(arr []float64) (i int) { //index max value in array
	max := arr[0]
	for j := 1; j < len(arr); j++ {
		if arr[j] > max {
			max = arr[j]
			i = j
		}
	}

	return i
}
