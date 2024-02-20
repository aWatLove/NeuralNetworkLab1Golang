package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
)

type NeuralNetwork struct {
	neurons [][]float64   // оутпуты на нейроне и сами нейроны
	w       [][][]float64 // веса [слой][нейрон][вес]

	LR    float64 // скорость обучения
	EPOCH int     // кол-во эпох обучения
}

//TODO: чтобы учитывалось смещение B biases + просчитывалось в forward() и при создании весов тоже. Можно сделать как Доп.Нейрон
//TODO: сделать обработку ошибок с тем что дата сет и кол-во входных нейронов совпадало
//TODO: сделать обработку ошибок с тем чтобы ожидаемые результаты с выходными как то правильно работали

func NewNeuralNetwork(neurons [][]float64, w [][][]float64, LR float64, EPOCH int) *NeuralNetwork {
	return &NeuralNetwork{neurons: neurons, w: w, LR: LR, EPOCH: EPOCH}
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

	neurons := createNeurons(inputData, 10, 8, 5, 3)
	w := createWeights(neurons)
	generateWeights(w)
	var nn = NeuralNetwork{neurons: neurons, w: w, LR: 0.1, EPOCH: 5000}

	fmt.Printf("NN:\n %v\n", nn.neurons)

	train(&nn, trainData, expRes)

	predict(&nn, trainData[0])
	out := imvia(nn.neurons[len(nn.neurons)-1])
	fmt.Println(nn.neurons[len(nn.neurons)-1])
	fmt.Println("result:", results[out])

	fmt.Printf("NN:\n %v\n", nn.neurons)

	predict(&nn, trainData[61])
	fmt.Printf("NN:\n %v\n", nn.neurons)
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

func saveWeights() { // сохранить мозги (веса) в файл // todo

}

func loadWeights() { // загрузить пресет мозгов (весов) // todo

}

func generateWeights(w [][][]float64) {
	//rand.Seed(1)
	for i := 0; i < len(w); i++ {
		for j := 0; j < len(w[i]); j++ {
			for k := 0; k < len(w[i][j]); k++ {
				// Инициализация весов методом "He"
				w[i][j][k] = rand.NormFloat64() * math.Sqrt(2.0/float64(len(w[i][j])))
			}
		}
	}
}

func activate(s float64) float64 { // функция активации
	//return math.Tanh(a * s) // гиперболический тангенс
	return (1 / (1 + math.Pow(math.E, -s)))
}

func forward(nn *NeuralNetwork) { // функция прямого распространения
	for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
		for i, _ := range nn.neurons[n+1] {
			var res float64 = 0
			for j, input := range nn.neurons[n] {
				res += input * nn.w[n][j][i] // i - к output нейрону, j - из input нейрона
			}
			nn.neurons[n+1][i] = activate(res)
		}
	}
}

func backProp(nn *NeuralNetwork, exp []float64) {
	forward(nn)
	lastLayer := len(nn.neurons) - 1
	m := createMatrixByNN(nn.neurons)

	// Вычисляем градиенты в последнем слое
	for i, n := range nn.neurons[lastLayer] {
		m[lastLayer][i] = n * (1 - n) * (n - exp[i])
	}

	// Обратное распространение ошибки через скрытые слои
	for i := lastLayer - 1; i > 0; i-- {
		for j := 0; j < len(nn.neurons[i]); j++ {
			var sum float64
			for k, elem := range nn.neurons[i+1] {
				sum += elem * nn.w[i][j][k]
			}
			m[i][j] = sum * nn.neurons[i][j] * (1 - nn.neurons[i][j])
			for k := range nn.neurons[i+1] {

				// TODO: пересчитать формулу дельт и deltaW. И пересчитать заново
				deltaW := -(nn.LR * m[i][j] * nn.neurons[i][j]) //todo: выделить в фуннцию // походу m[i][j] = 0 и оно не высчитывается()
				nn.w[i][j][k] += deltaW

			}
		}
	}
}

func train(nn *NeuralNetwork, data [][]float64, exp [][]float64) {
	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам
		for d := 0; d < len(data); d++ { // цикл по дата сету
			for i := 0; i < len(data[d]); i++ {
				nn.neurons[0][i] = data[d][i] // присваиваем входным нейронам данные из дата-сета
			}
			backProp(nn, exp[d]) // вычисляем ошибку и корректируем веса
		}
		accuracy, errExp := evaluate(nn, data, exp)
		fmt.Printf("Epoch: %d, Accuracy: %.6f%%, ResExp: %.6f\n", e+1, accuracy*100, errExp) // todo:функцию ошибки посчитать
	}
}

func predict(nn *NeuralNetwork, data []float64) []float64 { // вычислить
	for i, _ := range nn.neurons[0] {
		nn.neurons[0][i] = data[i]
	}
	forward(nn)

	ll := len(nn.neurons) - 1 // last layer
	outputs := make([]float64, len(nn.neurons[ll]))
	for i, _ := range nn.neurons[ll] {
		outputs[i] = nn.neurons[ll][i]
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
