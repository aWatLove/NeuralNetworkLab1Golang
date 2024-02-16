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
	fmt.Println(trainData, expRes)

	neurons := createNeurons(8, 5, 5, 3)
	w := createWeights(neurons)
	generateWeights(w)
	var nn = NeuralNetwork{neurons: neurons, w: w, LR: 0.7, EPOCH: 5000}

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
	// Проходим по всем записям и преобразуем числовые значения в формат float64
	var data [][]float64
	for _, row := range records[1:] {
		var vector []float64
		ind := dictAdd(row[0])
		ind *= 0.01
		vector = append(vector, ind)
		for _, value := range row[1 : len(row)-1] { // Начинаем с 2-го элемента, так как первые два это строковые значения
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			floatValue *= 0.001
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
	trainData = data

	return trainData, expResults, nil
}

// загрузка тестового data set'а из файла
// надо сделать так чтобы он вытаскивал всю csv. А потом уже отделять данные входные и выходыне.
func loadTestData() ([][]float64, error) {

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

func saveWeights() { // сохранить мозги (веса) в файл

}

func loadWeights() { // загрузить пресет мозгов (весов)

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

//func generateWeights(w [][][]float64) { // сгенерировать рандомно веса
//	rand.Seed(1)
//	for i := 0; i < len(w); i++ {
//		for j := 0; j < len(w[i]); j++ {
//			for k := 0; k < len(w[i][j]); k++ {
//				w[i][j][k] = float64(float64(rand.Float64() * 0.01))
//			}
//		}
//	}
//}

func activate(s float64) float64 { // функция активации
	//return math.Tanh(a * s) // гиперболический тангенс
	return (1 / (1 + math.Pow(math.E, -s)))
}

/*
func forward(inputNeurons []float64, outputNeurons []float64, w [][]float64) { // функция прямого распространения

	for i, _ := range outputNeurons {
		var res float64 = 0
		for j, input := range inputNeurons {
			res += input * w[j][i] // i - к output нейрону, j - из input нейрона
		}
		outputNeurons[i] = activate(res)
	}
}
*/

func forward(nn *NeuralNetwork) { // функция прямого распространения
	for n := 0; n < len(nn.w); n++ { // цикл по слоям нейронов
		//forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
		for i, _ := range nn.neurons[n+1] {
			var res float64 = 0
			for j, input := range nn.neurons[n] {
				res += input * nn.w[n][j][i] // i - к output нейрону, j - из input нейрона
			}
			nn.neurons[n+1][i] = activate(res)
		}
		//fmt.Println(nn.neurons[n])

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
				deltaW := -(nn.LR * m[i][j] * nn.neurons[i][j])
				nn.w[i][j][k] += deltaW
			}
		}
	}
}

//func backProp(nn NeuralNetwork, exp []float64) { // Функция обратного распространения ошибки
//	// должен тут вызывать forward() и что-то получать и дальше считать цену ошибки
//	forward(nn)
//	m := createMatrixByNN(nn.neurons)
//	lastLayer := len(nn.neurons) - 1
//	for i, n := range nn.neurons[lastLayer] {
//		//m[lastLayer][i] = (1 / (math.Pow(math.Cosh(n), 2))) * (n - exp[i]) // n - полученное значение, exp - ожидаемое значение // производная гиперболического тангенс
//		m[lastLayer][i] = n * (1 - n) * (n - exp[i]) // n - полученное значение, exp - ожидаемое значение
//	}
//	//fmt.Println("LastLayer m", m[lastLayer])
//	//fmt.Println(len(m))
//	//fmt.Println(len(nn.neurons))
//
//	for i := lastLayer - 1; i > 0; i-- { // neuron layer
//		for j := 0; j < len(nn.neurons[i]); j++ { // some neuron
//			var sum float64
//			for k, elem := range nn.neurons[i+1] { // recursive sum for this neuron
//				sum += elem * nn.w[i][j][k]
//			}
//			//fmt.Println("COSH", math.Cosh(nn.neurons[i][j]))
//			//fmt.Println("COSH in square", (math.Pow(math.Cosh(nn.neurons[i][j]), 2)))
//			//fmt.Println("P", (1 / (math.Pow(math.Cosh(nn.neurons[i][j]), 2))))
//			//fmt.Println("sum", sum)
//			//fmt.Println("m[i][j]", sum*(1/(math.Pow(math.Cosh(nn.neurons[i][j]), 2))))
//
//			//m[i][j] = sum * (1 / (math.Pow(math.Cosh(nn.neurons[i][j]), 2))) // нашли б // производная tanh
//			m[i][j] = sum * nn.neurons[i][j] * (1 - nn.neurons[i][j]) // нашли б
//
//			for k, _ := range nn.neurons[i+1] { // коррекция весов //todo: пересмотреть этот цикл
//				//fmt.Println(nn.LR)
//				//fmt.Println("M: ", m[i][j])
//				//fmt.Println("В нейроне", nn.neurons[i][j])
//				deltaW := -(nn.LR * m[i][j] * nn.neurons[i][j]) // dW - коррекция весов // TODO: здесь трабл!!!!
//				//fmt.Println("deltaW: ", deltaW)
//
//				nn.w[i][j][k] = nn.w[i][j][k] - deltaW // изменяем веса
//			}
//		}
//	}
//fmt.Println("MMMMMMMMMMMMMMMMMMMMMMMMMMMMMMm", m)
//fmt.Println(" ")
//for i := len(nn.neurons) - 2; i >= 0; i-- {
//	for j := 0; j < len(nn.neurons[i]); j++ {
//
//		for k, _ := range nn.neurons[i+1] { // коррекция весов //todo: пересмотреть этот цикл
//			//fmt.Println(nn.LR)
//			//fmt.Println("M: ", m[i][j])
//			//fmt.Println("В нейроне", nn.neurons[i][j])
//			deltaW := -(nn.LR * m[i][j] * nn.neurons[i][j]) // dW - коррекция весов // TODO: здесь трабл!!!!
//			//fmt.Println("deltaW: ", deltaW)
//
//			nn.w[i][j][k] = nn.w[i][j][k] - deltaW // изменяем веса
//		}
//	}
//}
//}

func train(nn *NeuralNetwork, data [][]float64, exp [][]float64) {
	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам
		for d := 0; d < len(data); d++ { // цикл по дата сету
			for i := 0; i < len(data[d]); i++ {
				nn.neurons[0][i] = data[d][i] // присваиваем входным нейронам данные из дата-сета
			}
			backProp(nn, exp[d]) // вычисляем ошибку и корректируем веса
		}
		accuracy := evaluate(nn, data, exp)
		fmt.Printf("Epoch: %d, Accuracy: %.2f%%\n", e+1, accuracy*100)
	}
}

//func train(nn NeuralNetwork, data [][]float64, exp [][]float64) {
//
//	for e := 0; e < nn.EPOCH; e++ { // цикл по эпохам
//
//		for d := 0; d < len(data); d++ { // цикл по дата сету // вынести в функцию
//			for i := 0; i < len(data[d]); i++ {
//				nn.neurons[0][i] = data[d][i] // присваиваем входным нейронам данные из дата-сета
//			}
//			//forward(nn)
//			//
//			//for n := 0; n < len(nn.w); n++ { // цикл по слоям весов нейронов
//			//	forward(nn.neurons[n], nn.neurons[n+1], nn.w[n]) // текущий слой, следующий слой, веса между этими слоями
//			//}
//			backProp(nn, exp[d]) // вычисляем ошибку и корректируем веса
//		}
//
//		accuracy := evaluate(nn, data, exp)
//		fmt.Printf("Epoch: %d, Accuracy: %.2f%%\n", e+1, accuracy*100)
//
//	}
//}

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

func evaluate(nn *NeuralNetwork, data [][]float64, exp [][]float64) (accuracy float64) {
	var count int
	for i, d := range data {
		outputs := predict(nn, d)
		o := imvia(outputs) // output
		t := imvia(exp[i])  // target - ожидаемое значение
		if o == t {
			count++
		}
	}

	accuracy = float64(count) / float64(len(data))
	return accuracy
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
