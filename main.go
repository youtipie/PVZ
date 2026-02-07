package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"os"
	"encoding/json"
	"sort"
)

// Структура для передачі даних у шаблони
type PageData struct {
	IsIndex bool
    Results       map[string]interface{}
    DefaultValues map[string]interface{}
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

	// Практика 3
	http.HandleFunc("/prac-3/task-1", prac3Task1)

	// Практика 4
	http.HandleFunc("/prac-4/task-1", prac4Task1)

	// Практика 5
    http.HandleFunc("/prac-5/task-1", prac5Task1)
	http.HandleFunc("/prac-5/data", prac5DataHandler) // API для отримання списку елементів

	// Практика 6
	http.HandleFunc("/prac-6/task-1", prac6Task1)

	log.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Допоміжна функція для рендеру темплейтів
func render(w http.ResponseWriter, tmplName string, data PageData, files ...string) {
	files = append([]string{"templates/base.html"}, files...)

	funcMap := template.FuncMap{
		"floatToStr": func(v interface{}) string {
			switch val := v.(type) {
			case float64:
				return fmt.Sprintf("%.2f", val)
			case int:
				return fmt.Sprintf("%d", val)
			default:
				return fmt.Sprintf("%v", val)
			}
		},
		"safeIndex": func(list interface{}, i int) interface{} {
			switch v := list.(type) {
			case []float64:
				if i >= 0 && i < len(v) {
					return v[i]
				}
			case []interface{}:
				if i >= 0 && i < len(v) {
					return v[i]
				}
			}
			return 0.0
		},
		"getResAtIndex": func(key string, i int, results map[string]interface{}) interface{} {
			if val, ok := results[key]; ok {
				if list, ok := val.([]float64); ok {
					if i >= 0 && i < len(list) {
						return fmt.Sprintf("%.2f", list[i])
					}
				}
			}
			return "-"
		},
		"getRes": func(key string, results map[string]interface{}) interface{} {
			if val, ok := results[key]; ok {
				return val
			}
			return "-"
		},
		"iterate": func(count int) []int {
			var items []int
			for i := 0; i < count; i++ {
				items = append(items, i)
			}
			return items
		},
		"add": func(a, b int) int {
			return a + b
		},
	}

	tmpl := template.New("base.html").Funcs(funcMap)

	var err error
	tmpl, err = tmpl.ParseFiles(files...)
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

func getFloatList(r *http.Request, key string) []float64 {
	var result []float64
	values := r.Form[key]
	for _, v := range values {
		v = strings.ReplaceAll(v, ",", ".")
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			result = append(result, f)
		} else {
			result = append(result, 0)
		}
	}
	return result
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
		data.Results = map[string]interface{}{
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
		data.Results = map[string]interface{}{
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
		data.Results = map[string]interface{}{
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

// Шлях, що обробляє перше завдання третьої практичної роботи
func prac3Task1(w http.ResponseWriter, r *http.Request) {
	data := PageData{IsIndex: false}

	if r.Method == http.MethodPost {
		// Отримання користувацього вводу
		Pc, err1 := getFloat(r, "Pc")
		q1, err2 := getFloat(r, "Q1")
		q2, err3 := getFloat(r, "Q2")
		B, err4 := getFloat(r, "B")

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			data.Error = "Bad values: check inputs"
			render(w, "prac_3_task_1", data, "templates/prac_3_task_1.html")
			return
		}

		// Якщо q2 більше, то це не має сенсу, сповіщаємо про помилку
		if q2 >= q1 {
			data.Error = "σ2 має бути менше за σ1."
			render(w, "prac_3_task_1", data, "templates/prac_3_task_1.html")
			return
		}

		// Функція для розрахунку прибутку
		// Використовуємо math.Erf для інтегрування нормального розподілу
		calculateProfit := func(sigma float64) float64 {
			// Межі інтегрування: Pc - 0.05*Pc до Pc + 0.05*Pc
			// Це симетричний інтервал навколо середнього (Pc).
			// Інтеграл від PDF нормального розподілу в межах [μ - δ, μ + δ] дорівнює erf(δ / (σ * sqrt(2)))
			delta := 0.05 * Pc
			qW := math.Erf(delta / (sigma * math.Sqrt(2)))

			// Розрахуємо прибуток (частка без небалансу)
			W_success := Pc * 24 * qW
			P_success := W_success * B

			// Розрахуємо штраф (частка з небалансом)
			W_imbalance := Pc * 24 * (1 - qW)
			Penalty := W_imbalance * B

			return P_success - Penalty
		}

		res1 := calculateProfit(q1)
		res2 := calculateProfit(q2)

		// Заносимо результати
		data.Results = map[string]interface{}{
			"res1": round(res1, 2),
			"res2": round(res2, 2),
			"q1":   q1,
			"q2":   q2,
		}
	}

	render(w, "prac_3_task_1", data, "templates/prac_3_task_1.html")
}

// Метод, що читає дані про кабеля з файлу
func getJek(index int, Tm float64) (float64, error) {
	file, err := os.Open("./instance/prac_4_cabels_data.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var data map[string][]float64
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return 0, err
	}

	var key string
	if Tm >= 1000 && Tm <= 3000 {
		key = "1000-3000"
	} else if Tm > 3000 && Tm <= 5000 {
		key = "3000-5000"
	} else if Tm > 5000 {
		key = "5000+"
	} else {
		return 0, fmt.Errorf("Tm out of range")
	}

	if vals, ok := data[key]; ok {
		if index >= 0 && index < len(vals) {
			return vals[index], nil
		}
	}
	return 0, fmt.Errorf("data not found for index %d", index)
}

// Метод, що "заокруглює" значення перерізу кабеля (шукає найближче стандартне)
func getCrossSection(value float64) float64 {
	crossSections := []float64{10, 16, 25, 35, 50, 70, 95, 120, 150, 185, 240}
	closest := crossSections[0]
	minDiff := math.Abs(value - closest)

	for _, s := range crossSections {
		diff := math.Abs(value - s)
		if diff < minDiff {
			minDiff = diff
			closest = s
		}
	}
	return closest
}

// Шлях, що обробляє четверту практичну роботу
func prac4Task1(w http.ResponseWriter, r *http.Request) {
	// Значення за замовчуванням
	defaultValues := map[string]interface{}{
		"Ik": 2500.0,
		"tf": 2.5,
		"Sm": 1300.0,
		"Tm": 4000.0,
		"Sk": 200.0,
	}
	data := PageData{
		IsIndex:       false,
		DefaultValues: defaultValues,
	}

	if r.Method == http.MethodPost {
		cabelStr := r.FormValue("cabel")
		cabel, errC := strconv.Atoi(cabelStr)
		Ik, err1 := getFloat(r, "Ik")
		tf, err2 := getFloat(r, "tf")
		Sm, err3 := getFloat(r, "Sm")
		Tm, err4 := getFloat(r, "Tm")
		Sk, err5 := getFloat(r, "Sk")

		if errC != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
			data.Error = "Bad values: check inputs"
			render(w, "prac_4_task_1", data, "templates/prac_4_task_1.html")
			return
		}

		// Оновлюємо значення за замовчуванням на введені користувачем
		data.DefaultValues = map[string]interface{}{
			"Ik": Ik,
			"tf": tf,
			"Sm": Sm,
			"Tm": Tm,
			"Sk": Sk,
		}

		// 1
		// Розрахунковий струм для нормального і післяаварійного режимів
		Im := (Sm / 2) / (math.Sqrt(3) * 10)
		Im_pa := 2 * Im

		// Отримуємо економічну густину струму
		jek, errJ := getJek(cabel, Tm)
		if errJ != nil {
			data.Error = "Cable data error: " + errJ.Error()
			render(w, "prac_4_task_1", data, "templates/prac_4_task_1.html")
			return
		}

		// Рахуємо економічний переріз
		sek := Im / jek
		// Шукаємо мінімальний переріз
		s_min := (Ik * math.Sqrt(tf)) / 92
		// На основі мінімального перерізу шукаємо кабель з потрібним перерізом
		s := getCrossSection(s_min)

		// 2
		// Рауємо опори елементів
		Xc := math.Pow(10.5, 2) / Sk
		Xt := (10.5 / 100) * (math.Pow(10.5, 2) / 6.3)
		// Сумарний опір
		Xe := Xc + Xt
		// Початкове діюче значення струму трифазного КЗ
		Ip0 := 10.5 / (math.Sqrt(3) * Xe)

		// 3
		// Сталі дані, передані з підстанції
		Rcn := 10.65
		Xcn := 24.02
		Rcmin := 34.88
		Xcmin := 65.68
		Uk_max := 11.1
		Uvn := 115.0
		Unn := 11.0
		Snomt := 6.3

		// Розрахуємо реактивний опір силового трансформатора
		Xt_tr := (Uk_max * math.Pow(Uvn, 2)) / (100 * Snomt)

		// Розрахуємо опори на шинах 10 кВ в нормальному та мінімальному режимах
		Rsh := Rcn
		Xsh := Xcn + Xt_tr
		Zsh := math.Sqrt(math.Pow(Rsh, 2) + math.Pow(Xsh, 2))

		Rshmin := Rcmin
		Xshmin := Xcmin + Xt_tr
		Zshmin := math.Sqrt(math.Pow(Rshmin, 2) + math.Pow(Xshmin, 2))

		// Розраховуємо струми трифазного та двофазного КЗ на шинах 10 кВ
		Ish_3 := (Uvn * 1000) / (math.Sqrt(3) * Zsh)
		Ish_2 := Ish_3 * math.Sqrt(3) / 2

		Ish_min_3 := (Uvn * 1000) / (math.Sqrt(3) * Zshmin)
		Ish_min_2 := Ish_min_3 * math.Sqrt(3) / 2

		// Розраховуємо коефіцієнт приведення
		kpr := math.Pow(Unn, 2) / math.Pow(Uvn, 2)

		// Розраховуємо опори на шинах 10 кВ в нормальному
		// та мінімальному режимах і заносимо їх в карту вставок
		Rshn := Rsh * kpr
		Xshn := Xsh * kpr
		Zshn := math.Sqrt(math.Pow(Rshn, 2) + math.Pow(Xshn, 2))

		Rshn_min := Rshmin * kpr
		Xshn_min := Xshmin * kpr
		Zshn_min := math.Sqrt(math.Pow(Rshn_min, 2) + math.Pow(Xshn_min, 2))

		// Розраховуємо дійсні струми трифазного та двофазного КЗ
		Ishn_3 := (Unn * 1000) / (math.Sqrt(3) * Zshn)
		Ishn_2 := Ishn_3 * math.Sqrt(3) / 2

		Ishn_min_3 := (Unn * 1000) / (math.Sqrt(3) * Zshn_min)
		Ishn_min_2 := Ishn_min_3 * math.Sqrt(3) / 2

		// Розрахунок струмів короткого замикання відхідних ліній 10 кВ
		R0 := 0.64
		X0 := 0.363
		// Знайдемо резистанси та реактанси відрізка з найбільшим опором
		Il := 0.2 + 0.35 + 0.2 + 0.6 + 2 + 2.55 + 3.37 + 3.1
		Rl := Il * R0
		Xl := Il * X0

		// Розрахуємо опори в нормальному та мінімальному режимах
		Ren := Rl + Rshn
		Xen := Xl + Xshn
		Zen := math.Sqrt(math.Pow(Ren, 2) + math.Pow(Xen, 2))

		Ren_min := Rl + Rshn_min
		Xen_min := Xl + Xshn_min
		Zen_min := math.Sqrt(math.Pow(Ren_min, 2) + math.Pow(Xen_min, 2))

		// Розрахуємо струми трифазного і двофазного КЗ
		Iln_3 := (Unn * 1000) / (math.Sqrt(3) * Zen)
		Iln_2 := Iln_3 * math.Sqrt(3) / 2

		Iln_min_3 := (Unn * 1000) / (math.Sqrt(3) * Zen_min)
		Iln_min_2 := Iln_min_3 * math.Sqrt(3) / 2

		// Заносимо усі результати у список
		data.Results = map[string]interface{}{
			"sek":        round(sek, 2),
			"s":          s,
			"Im":         round(Im, 2),
			"Im_pa":      round(Im_pa, 2),
			"Ip0":        round(Ip0, 2),
			"Ish_3":      round(Ish_3, 2),
			"Ish_2":      round(Ish_2, 2),
			"Ish_min_3":  round(Ish_min_3, 2),
			"Ish_min_2":  round(Ish_min_2, 2),
			"Ishn_3":     round(Ishn_3, 2),
			"Ishn_2":     round(Ishn_2, 2),
			"Ishn_min_3": round(Ishn_min_3, 2),
			"Ishn_min_2": round(Ishn_min_2, 2),
			"Iln_3":      round(Iln_3, 2),
			"Iln_2":      round(Iln_2, 2),
			"Iln_min_3":  round(Iln_min_3, 2),
			"Iln_min_2":  round(Iln_min_2, 2),
		}
	}
	render(w, "prac_4_task_1", data, "templates/prac_4_task_1.html")
}

// Метод, що читає дані елементів ЕПС з JSON файлу
func getPrac5Data() (map[string][]float64, error) {
	file, err := os.Open("./instance/prac_5_data.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string][]float64
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// API Handler для отримання списку елементів
func prac5DataHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getPrac5Data()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

// Шлях, що обробляє п'яту практичну роботу
func prac5Task1(w http.ResponseWriter, r *http.Request) {
	data := PageData{IsIndex: false}

	if r.Method == http.MethodPost {
		// Отримання користувацього вводу
		r.ParseForm()
		quantitiesStr := r.Form["quantity[]"]
		elements := r.Form["element[]"]
		Zpera, err1 := getFloat(r, "Zpera")
		Zperp, err2 := getFloat(r, "Zperp")

		if err1 != nil || err2 != nil {
			data.Error = "Bad values: check inputs"
			render(w, "prac_5_task_1", data, "templates/prac_5_task_1.html")
			return
		}

		pracData, err := getPrac5Data()
		if err != nil {
			data.Error = "Error reading data file"
			render(w, "prac_5_task_1", data, "templates/prac_5_task_1.html")
			return
		}

		var woc_sum, tvoc_num, max_t_plan float64

		for i, el := range elements {
			if i >= len(quantitiesStr) {
				break
			}
			q, err := strconv.Atoi(quantitiesStr[i])
			if err != nil {
				continue
			}
			quantity := float64(q)

			if props, ok := pracData[el]; ok && len(props) >= 3 {
				omega := props[0]
				tv := props[1]
				tp := props[2]

				woc_sum += quantity * omega
				tvoc_num += quantity * omega * tv

				if tp > max_t_plan {
					max_t_plan = tp
				}
			}
		}

		// Розрахунки
		woc := woc_sum
		var tvoc float64
		if woc > 0 {
			tvoc = tvoc_num / woc
		}

		// Коефіцієнт аварійного простою одноколової системи
		kaoc := (woc * tvoc) / 8760
		// Коефіцієнт планового простою одноколової системи
		kpoc := (1.2 * max_t_plan) / 8760
		// Частота відмов одночасно двох кіл двоколової системи
		wdk := 2 * woc * (kaoc + kpoc)
		// Частота відмов двоколової системи з урахуванням секційного вимикача
		wdc := wdk + 0.02
		// Коефіцієнт надійності
		var koef float64
		if wdc > 0 {
			koef = woc / wdc
		}

		// Пункт 2
		w := 0.01
		tv := 0.045 // 45 * 10^-3
		Pm := 5120.0 // 5.12 * 10^3
		Tm := 6451.0
		kp := 0.004 // 4 * 10^-3

		M_1 := w * tv * Pm * Tm
		M_2 := kp * Pm * Tm
		M := Zpera*M_1 + Zperp*M_2

		data.Results = map[string]interface{}{
			"woc":  round(woc, 4),
			"wdc":  round(wdc, 4),
			"koef": koef,
			"M":    round(M, 0),
		}
	}

	render(w, "prac_5_task_1", data, "templates/prac_5_task_1.html")
}

// Метод, щоб знайти найближчі межі до числа у списку(Використовуємо при пошуці Кв)
func findNearestNeighbors(lst []int, target int) (int, int) {
	sort.Ints(lst)
	lower := -1
	higher := -1
	for _, v := range lst {
		if v < target {
			lower = v
		} else if v > target {
			higher = v
			break
		}
	}
	return lower, higher
}

// Метод для пошуку значення розрахункових коефіцієнтів Кр
// для мереж живлення напругою до 1000 В (Т0 = 10 хв.)
func getKp1(ne int, groupUseCoff float64) (float64, error) {
	file, err := os.Open("./instance/prac_6_data_1.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return 0, err
	}

	// Визначити найближчий коефіцієнтний ключ.
    // Ключі мають вигляд «0,1», «0,15» тощо.
    // Потрібно знайти максимальний ключ <= groupUseCoff.
    // Припустимо, що рядок «1» існує, щоб отримати ключі.
	firstRow := data["1"]
	var coeffs []float64
	for k := range firstRow {
		f, _ := strconv.ParseFloat(k, 64)
		coeffs = append(coeffs, f)
	}
	sort.Float64s(coeffs)

	closestCoeff := coeffs[0]
	for _, c := range coeffs {
		if c <= groupUseCoff {
			closestCoeff = c
		} else {
			break
		}
	}
	coeffKey := fmt.Sprintf("%g", closestCoeff)

	if row, ok := data[strconv.Itoa(ne)]; ok {
		if val, ok := row[coeffKey]; ok {
			return val, nil
		}
	}

	// Інтерполяція
	var intKeys []int
	for k := range data {
		i, _ := strconv.Atoi(k)
		intKeys = append(intKeys, i)
	}
	lower, higher := findNearestNeighbors(intKeys, ne)

	if lower != -1 && higher != -1 {
		rowLower := data[strconv.Itoa(lower)]
		rowHigher := data[strconv.Itoa(higher)]

		valLower := rowLower[coeffKey]
		valHigher := rowHigher[coeffKey]

		slope := (valHigher - valLower) / float64(higher - lower)
		value := valLower + slope*float64(ne-lower)
		return value, nil
	}

	return 0, fmt.Errorf("Kp1 lookup failed for ne=%d", ne)
}

func getKp2(ne int, groupUseCoff float64) (float64, error) {
	file, err := os.Open("./instance/prac_6_data_2.json")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var data map[string]map[string]float64
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return 0, err
	}

	var targetRow map[string]float64
	for k, row := range data {
		parts := strings.Split(k, ";")
		if len(parts) == 2 {
			minV, _ := strconv.Atoi(parts[0])
			maxV, _ := strconv.ParseInt(parts[1], 10, 64)
			if ne >= minV && int64(ne) <= maxV {
				targetRow = row
				break
			}
		}
	}

	if targetRow == nil {
		return 0, fmt.Errorf("range not found for ne=%d", ne)
	}

	var coeffs []float64
	for k := range targetRow {
		f, _ := strconv.ParseFloat(k, 64)
		coeffs = append(coeffs, f)
	}
	sort.Float64s(coeffs)

	closestCoeff := coeffs[0]
	for _, c := range coeffs {
		if c <= groupUseCoff {
			closestCoeff = c
		} else {
			break
		}
	}
	coeffKey := fmt.Sprintf("%g", closestCoeff)

	if val, ok := targetRow[coeffKey]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("Kp2 coeff not found")
}

func prac6Task1(w http.ResponseWriter, r *http.Request) {
	// Отримуємо значення по змовчуванню для таблиці (Значення з контрольного прикладу)
	file, err := os.Open("./instance/prac_6_table_default_data.json")
	if err != nil {
		http.Error(w, "Failed to load default data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var defaultValues map[string]interface{}
	if err := json.NewDecoder(file).Decode(&defaultValues); err != nil {
		http.Error(w, "Failed to parse default data", http.StatusInternalServerError)
		return
	}

	data := PageData{IsIndex: false, DefaultValues: defaultValues, Results: make(map[string]interface{})}

	if r.Method == http.MethodPost {
		r.ParseForm()

		// Допоміжна функція для парсингу списків
		parseList := func(key string) []float64 {
			return getFloatList(r, key)
		}

		// Отримуємо користувацький ввід для ЕП першого ШР
		nu := parseList("nu[]")
		cos := parseList("cos[]")
		Uh := parseList("Uh[]")
		n := parseList("n[]")
		Ph := parseList("Ph[]")
		KB := parseList("KB[]")
		tg := parseList("tg[]")

		count := len(nu)

		// Масиви для результатів
		var nPh, Ip, nPhKB, nPhKBtg, nPhSquare []float64
		var sum_nPh, sum_nPhKB, sum_nPhKBtg, sum_nPhSquare float64
		var sum_n float64

		// Шукаємо розрахункові струми на І рівні електропостачання
		for i := 0; i < count; i++ {
			// n * Ph
			val_nPh := n[i] * Ph[i]
			nPh = append(nPh, val_nPh)
			sum_nPh += val_nPh

			// Знаходимо розрахунковий струм кожного ЕП
			val_Ip := val_nPh / (math.Sqrt(3) * Uh[i] * cos[i] * nu[i])
			Ip = append(Ip, round(val_Ip, 2))

			// n * Ph * KB
			val_nPhKB := val_nPh * KB[i]
			nPhKB = append(nPhKB, round(val_nPhKB, 2))
			sum_nPhKB += val_nPhKB

			// n * Ph * KB * tg
			val_nPhKBtg := val_nPhKB * tg[i]
			nPhKBtg = append(nPhKBtg, round(val_nPhKBtg, 2))
			sum_nPhKBtg += val_nPhKBtg

			// n * Ph^2
			val_nPhSquare := n[i] * math.Pow(Ph[i], 2)
			nPhSquare = append(nPhSquare, round(val_nPhSquare, 2))
			sum_nPhSquare += val_nPhSquare

			sum_n += n[i]
		}

		// Знаходимо груповий коефіцієнт використання
		groupUseCoff := 0.0
		if sum_nPh > 0 {
			groupUseCoff = sum_nPhKB / sum_nPh
		}

		// Знаходимо ефективну кількість ЕП
		ne := 0
		if sum_nPhSquare > 0 {
			ne = int(math.Ceil(math.Pow(sum_nPh, 2) / sum_nPhSquare))
		}

		// Знаходимо розрахунковий коефіцієнт активної потужності по таблиці 3.3
		// за допомогою методу, описаного раніше в utils.py
		Kp, err := getKp1(ne, groupUseCoff)
		if err != nil {
			Kp = 0
			fmt.Println("Error finding Kp1:", err)
		}

		// Знаходимо розрахункове активне навантаження
		Pp := Kp * sum_nPhKB
		// Знаходимо розрахункове реактивне навантаження
		Qp := Kp * sum_nPhKBtg
		// Знаходимо повну потужність
		Sp := math.Sqrt(math.Pow(Pp, 2) + math.Pow(Qp, 2))

		// Середня напруга для розрахунку загального струму
		var sumUh float64
		for _, v := range Uh {
			sumUh += v
		}
		meanUh := 0.38 // Значення за замовчуванням
		if len(Uh) > 0 {
			meanUh = sumUh / float64(len(Uh))
		}

		// Знаходимо розрахунковий груповий струм ШР1
		Ip_total := Pp / meanUh

		// Отримуємо користувацький ввід для крупних ЕП
		nu_big := parseList("nu_big[]")
		cos_big := parseList("cos_big[]")
		Uh_big := parseList("Uh_big[]")
		n_big := parseList("n_big[]")
		Ph_big := parseList("Ph_big[]")
		KB_big := parseList("KB_big[]")
		tg_big := parseList("tg_big[]")

		var nPh_big, Ip_big, nPhKB_big, nPhKBtg_big, nPhSquare_big []float64

		// Розрахунки ті ж самі, що і для звичайних ЕП
		for i := 0; i < len(nu_big); i++ {
			v_nPh := n_big[i] * Ph_big[i]
			nPh_big = append(nPh_big, v_nPh)

			v_nPhKB := v_nPh * KB_big[i]
			nPhKB_big = append(nPhKB_big, round(v_nPhKB, 2))

			// Для другого ЕП відсутнє значення коефіцієнту реактивної потужності відсутнє,
			// тому замість нього пишемо 0 (Це ніяк не вплине на розрахунки, зроблено тільки для зручності
			v_tg := 0.0
			if i < len(tg_big) {
				v_tg = tg_big[i]
			}
			v_nPhKBtg := v_nPhKB * v_tg
			nPhKBtg_big = append(nPhKBtg_big, round(v_nPhKBtg, 2))

			v_nPhSq := n_big[i] * math.Pow(Ph_big[i], 2)
			nPhSquare_big = append(nPhSquare_big, round(v_nPhSq, 2))

			v_Ip := v_nPh / (math.Sqrt(3) * Uh_big[i] * cos_big[i] * nu_big[i])
			Ip_big = append(Ip_big, round(v_Ip, 2))
		}

		// Отримуємо користувацький ввід загального навантаження цеху
		n_all, _ := getFloat(r, "n")
		nPh_all, _ := getFloat(r, "nPh")
		nPhKB_all, _ := getFloat(r, "nPhKB")
		nPhKBtg_all, _ := getFloat(r, "nPhKBtg")
		nPhSquare_all, _ := getFloat(r, "nPh_square")

		// Знаходимо коефіцієнти використання цеху в цілому
		groupUseCoffAll := 0.0
		if nPh_all > 0 {
			groupUseCoffAll = nPhKB_all / nPh_all
		}

		// Знаходимо ефективну кількість ЕП цеху в цілому
		neAll := 0
		if nPhSquare_all > 0 {
			neAll = int(math.Round(math.Pow(nPh_all, 2) / nPhSquare_all))
		}

		// Знаходимо розрахунковий коефіцієнт активної потужності по таблиці 3.4
		// за допомогою методу, описаного раніше в utils.py
		KpAll, _ := getKp2(neAll, groupUseCoffAll)

		// Знаходимо розрахункове активне навантаження на шинах 0,38 кВ ТП
		PpAll := KpAll * nPhKB_all
		// Знаходимо розрахункове реактивне навантаження на шинах 0,38 кВ ТП
		QpAll := KpAll * nPhKBtg_all
		// Знаходимо повну потужність на шинах 0,38 кВ ТП
		SpAll := math.Sqrt(math.Pow(PpAll, 2) + math.Pow(QpAll, 2))

		// Середня напруга (використовуємо напругу крупних ЕП як базу 0.38)
		var sumUhBig float64
		for _, v := range Uh_big {
			sumUhBig += v
		}
		meanUhBig := 0.38
		if len(Uh_big) > 0 {
			meanUhBig = sumUhBig / float64(len(Uh_big))
		}

		// Знаходимо розрахунковий груповий струм на шинах 0,38 кВ ТП
		IpAll := PpAll / meanUhBig

		// Заносимо усі результати у список
		results := map[string]interface{}{
			"nPh_list": nPh, "Ip_list": Ip, "nPhKB_list": nPhKB, "nPhKBtg_list": nPhKBtg, "nPh_square_list": nPhSquare,
			"group_use_coff": round(groupUseCoff, 1),
			"ne":             ne, "Kp": round(Kp, 2), "Pp": round(Pp, 2), "Qp": round(Qp, 2), "Sp": round(Sp, 2), "Ip": round(Ip_total, 2),
			"N": int(sum_n), "nPh_sum": int(sum_nPh), "nPhKB_sum": round(sum_nPhKB, 2), "nPhKBtg_sum": round(sum_nPhKBtg, 2),
			"nPh_square_sum":    round(sum_nPhSquare, 2),
			"nPh_big_list":      nPh_big, "Ip_big_list": Ip_big, "nPhKB_big_list": nPhKB_big,
			"nPhKBtg_big_list":  nPhKBtg_big, "nPh_square_big_list": nPhSquare_big,
			"group_use_coff_all": round(groupUseCoffAll, 2), "ne_all": neAll, "Kp_all": round(KpAll, 2),
			"Pp_all": round(PpAll, 2), "Qp_all": round(QpAll, 2), "Sp_all": round(SpAll, 2), "Ip_all": round(IpAll, 2),
		}
		data.Results = results

		// Також створюємо список, який позначає користувацьки ввід
		// Це створено для того, щоб після розрахунків, значення введені користувачем, лишились
		userValues := make(map[string]interface{})

		normalMap := make(map[string]interface{})
		normalMap["naming"] = defaultValues["normal"].(map[string]interface{})["naming"]
		normalMap["nu[]"] = nu
		normalMap["cos[]"] = cos
		normalMap["Uh[]"] = Uh
		normalMap["n[]"] = n
		normalMap["Ph[]"] = Ph
		normalMap["KB[]"] = KB
		normalMap["tg[]"] = tg
		userValues["normal"] = normalMap

		bigMap := make(map[string]interface{})
		bigMap["naming"] = defaultValues["big"].(map[string]interface{})["naming"]
		bigMap["nu[]"] = nu_big
		bigMap["cos[]"] = cos_big
		bigMap["Uh[]"] = Uh_big
		bigMap["n[]"] = n_big
		bigMap["Ph[]"] = Ph_big
		bigMap["KB[]"] = KB_big
		bigMap["tg[]"] = tg_big
		userValues["big"] = bigMap

		allMap := make(map[string]interface{})
		allMap["n"] = n_all
		allMap["nPh"] = nPh_all
		allMap["nPhKB"] = nPhKB_all
		allMap["nPhKBtg"] = nPhKBtg_all
		allMap["nPh_square"] = nPhSquare_all
		userValues["all"] = allMap

		data.DefaultValues = userValues
	}

	render(w, "prac_6_task_1", data, "templates/prac_6_task_1.html")
}