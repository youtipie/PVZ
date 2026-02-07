package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Структура для передачі даних у шаблони
type PageData struct {
	IsIndex bool
	Results map[string]float64
	Error   string
}

func main() {
	// Обробка статичних файлів (js/css), якщо вони є
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Шлях, що обробляє головну сторінку, з якої можна потрапити на усі веб калькулятори курсу
	http.HandleFunc("/", index)

	// Практика 1
	http.HandleFunc("/prac-1/task-1", prac1Task1)
	http.HandleFunc("/prac-1/task-2", prac1Task2)

	// Практика 2
	http.HandleFunc("/prac-2/task-1", prac2Task1)

	log.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Допоміжна функція для рендеру темплейтів
func render(w http.ResponseWriter, tmplName string, data PageData, files ...string) {
	files = append([]string{"templates/base.html"}, files...)
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Render error: %v", err), http.StatusInternalServerError)
	}
}

// Допоміжна функція для парсингу чисел з плаваючею точкою
func getFloat(r *http.Request, key string) (float64, error) {
	val := r.FormValue(key)
	val = strings.ReplaceAll(val, ",", ".")
	if val == "" {
		return 0, fmt.Errorf("empty value for %s", key)
	}
	return strconv.ParseFloat(val, 64)
}

// Допоміжна функція для заокруглення чисел
func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// Шлях, що обробляє головну сторінку, з якої можна потрапити на усі веб калькулятори курсу
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{IsIndex: true}
	render(w, "index", data, "templates/index.html")
}

// Шлях, що обробляє перше завдання першої практичної роботи
// Оскільки на сторінці присутня форма, для користувацього вводу,
// тому додаємо можливість обробки POST запиту
func prac1Task1(w http.ResponseWriter, r *http.Request) {
	data := PageData{IsIndex: false}

	// Перевіряємо який запит було здійснено
	// Якщо POST, то на цю ж сторінку передаємо результати обрахунків
	if r.Method == http.MethodPost {
		// Отримання користувацього вводу
		// Також заміняємо ',' на '.' для коректного переведення строки у float
		// (реалізовано у функції getFloat)
		Hp, err1 := getFloat(r, "Hp")
		Cp, err2 := getFloat(r, "Cp")
		Sp, err3 := getFloat(r, "Sp")
		Np, err4 := getFloat(r, "Np")
		Op, err5 := getFloat(r, "Op")
		Wp, err6 := getFloat(r, "Wp")
		Ap, err7 := getFloat(r, "Ap")

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
			// Якщо в ході обрахунків виникає помилка, то виводимо помилку
			data.Error = "Bad values: check inputs"
			render(w, "prac_1_task_1", data, "templates/prac_1_task_1.html")
			return
		}

		// Обчислення результатів

		// Обчислюємо коефіцієнт переходу від робочої до сухої маси та
		// коефіцієнт переходу від робочої до горючої маси
		Kpc := 100 / (100 - Wp)
		Kpg := 100 / (100 - Wp - Ap)

		// Обчислюємо нижчу теплоту згоряння для робочої, сухої та горючої маси
		Qph := (339*Cp + 1030*Hp - 108.8*(Op-Sp) - 25*Wp) / 1000
		Qch := (Qph + 0.025*Wp) * 100 / (100 - Wp)
		Qgh := (Qph + 0.025*Wp) * 100 / (100 - Wp - Ap)

		// Обчислюємо склад сухої маси палива
		Hc := Hp * Kpc
		Cc := Cp * Kpc
		Sc := Sp * Kpc
		Nc := Np * Kpc
		Oc := Op * Kpc
		Ac := Ap * Kpc

		// Обчислюємо склад горючої маси палива
		Hg := Hp * Kpg
		Cg := Cp * Kpg
		Sg := Sp * Kpg
		Ng := Np * Kpg
		Og := Op * Kpg

		// Заносимо результати у словник (map) та округлюємо їх
		data.Results = map[string]float64{
			"Kpc": round(Kpc, 2),
			"Kpg": round(Kpg, 2),
			"Qph": round(Qph, 4),
			"Qch": round(Qch, 4),
			"Qgh": round(Qgh, 4),
			"Hc":  round(Hc, 2),
			"Cc":  round(Cc, 2),
			"Sc":  round(Sc, 2),
			"Nc":  round(Nc, 2),
			"Oc":  round(Oc, 2),
			"Ac":  round(Ac, 2),
			"Hg":  round(Hg, 2),
			"Cg":  round(Cg, 2),
			"Sg":  round(Sg, 2),
			"Ng":  round(Ng, 2),
			"Og":  round(Og, 2),
		}
	}

	// Рендеримо сторінку разом з результатами обрахунків (або без них для GET)
	render(w, "prac_1_task_1", data, "templates/prac_1_task_1.html")
}

// Шлях, що обробляє друге завдання першої практичної роботи
func prac1Task2(w http.ResponseWriter, r *http.Request) {
	data := PageData{IsIndex: false}

	if r.Method == http.MethodPost {
		// Код для другого завдання схожий:
		// Отримання користувацього вводу
		Hg, err1 := getFloat(r, "Hg")
		Cg, err2 := getFloat(r, "Cg")
		Sg, err3 := getFloat(r, "Sg")
		Vg, err4 := getFloat(r, "Vg")
		Og, err5 := getFloat(r, "Og")
		Wg, err6 := getFloat(r, "Wg")
		Ag, err7 := getFloat(r, "Ag")
		Qi, err8 := getFloat(r, "Qi")

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || err8 != nil {
			data.Error = "Bad values: check inputs"
			render(w, "prac_1_task_2", data, "templates/prac_1_task_2.html")
			return
		}

		// Обчислення результатів

		// Обчислюємо склад робочої маси мазуту
		Hp := Hg * (100 - Wg - Ag) / 100
		Cp := Cg * (100 - Wg - Ag) / 100
		Sp := Sg * (100 - Wg - Ag) / 100
		Op := Og * (100 - Wg - Ag) / 100
		Ap := Ag * (100 - Wg) / 100
		Vp := Vg * (100 - Wg) / 100

		// Обчислюємо нижчу теплоту згоряння мазуту на робочу масу
		Qri := Qi*(100-Wg-Ag)/100 - 0.025*Wg

		// Заносимо результати у словник та округлюємо їх
		data.Results = map[string]float64{
			"Hp":  round(Hp, 2),
			"Cp":  round(Cp, 2),
			"Sp":  round(Sp, 2),
			"Op":  round(Op, 2),
			"Ap":  round(Ap, 2),
			"Vp":  round(Vp, 2),
			"Qri": round(Qri, 4),
		}
	}

	// Рендеримо сторінку разом з результатами обрахунків
	render(w, "prac_1_task_2", data, "templates/prac_1_task_2.html")
}

// Шлях, що обробляє перше завдання другої практичної роботи
func prac2Task1(w http.ResponseWriter, r *http.Request) {
	data := PageData{IsIndex: false}

	if r.Method == http.MethodPost {
		// Отримання користувацього вводу
		coal, err1 := getFloat(r, "coal")
		oil, err2 := getFloat(r, "oil")
		// gas := getFloat(r, "gas") // Газ запитується, але не впливає на викиди твердих частинок

		// Також отримуємо константи, які може задати користувач
		Ap_coal, err3 := getFloat(r, "Ap")
		Qpi_coal, err4 := getFloat(r, "Qpi")
		Qgi_oil, err5 := getFloat(r, "Qgi_oil")
		Wp_oil, err6 := getFloat(r, "Wp_oil")
		Gvun, err7 := getFloat(r, "Gvun")
		nzu, err8 := getFloat(r, "nzu")

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil || err8 != nil {
			data.Error = "Bad values: check inputs"
			render(w, "prac_2_task_1", data, "templates/prac_2_task_1.html")
			return
		}

		// Значення частки леткої золи для вугілля та мазуту
		avun_coal := 0.8
		avun_oil := 1.0

		// Шукаємо нижчу теплоту згоряння робочї маси для мазуту
		Qri_oil := Qgi_oil*(100-Wp_oil-0.15)/100 - 0.025*Wp_oil

		// Обчислюємо показник емісії твердих частинок при спалюванні вугілля
		ktv_coal := math.Pow(10, 6) / Qpi_coal * avun_coal * Ap_coal / (100 - Gvun) * (1 - nzu)
		Etv_coal := math.Pow(10, -6) * ktv_coal * Qpi_coal * coal

		// Обчислюємо показник емісії твердих частинок при спалюванні мазуту
		ktv_oil := math.Pow(10, 6) / Qri_oil * avun_oil * 0.15 / 100 * (1 - nzu)
		Etv_oil := math.Pow(10, -6) * ktv_oil * Qri_oil * oil

		// Для газу = 0
		ktv_gas := 0.0
		Etv_gas := 0.0

		// Заносимо результати
		data.Results = map[string]float64{
			"ktv_coal": round(ktv_coal, 2),
			"Etv_coal": round(Etv_coal, 2),
			"ktv_oil":  round(ktv_oil, 2),
			"Etv_oil":  round(Etv_oil, 2),
			"ktv_gas":  ktv_gas,
			"Etv_gas":  Etv_gas,
		}
	}

	render(w, "prac_2_task_1", data, "templates/prac_2_task_1.html")
}
