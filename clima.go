package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Clima struct {
	Temperatura float64 `json:"temperatura"`
	Humedad     int     `json:"humedad"`
	Descripcion string  `json:"descripcion"`
	Viento      float64 `json:"viento"`
}

func obtenerClima() (Clima, error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=-0.181577&lon=-78.486571&appid=3a097cd255f919c4993f29a0a94d7285&units=metric&lang=es")
	if err != nil {
		return Clima{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return Clima{}, fmt.Errorf("error al consultar openweathermap. Status %d: Body %q", resp.StatusCode, string(b))
	}

	type owmResp struct {
		Main struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		} `json:"main"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`
	}

	var owm owmResp
	if err := json.NewDecoder(resp.Body).Decode(&owm); err != nil {
		return Clima{}, err
	}

	return Clima{
		Temperatura: owm.Main.Temp,
		Humedad:     owm.Main.Humidity,
		Descripcion: owm.Weather[0].Description,
		Viento:      owm.Wind.Speed,
	}, nil
}
