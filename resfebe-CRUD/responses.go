/*
 * Copyright (C) 2019  Murat Koptur
 *
 * Contact: mkoptur3@gmail.com
 *
 * Last edit: 3/30/19 10:21 PM
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"net/http"
)

// code: response statusu,400,500,404 vb.
// payload requestin en altında, headerın aşağısında ekstra verdiklerin. http request.

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	//object oriented değil interface var, payload interface, bunu jsona çeviriyorum.
	// header'ı json döneceğim için jsona set ediyorum
	// response'un headerına yazıyorsun. body'e de response'u json olarak yazdırıyorsun. payload=body yani

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
	// _,_ iki şey döndürecek write şeyi. ilgilendirmediği için _,_ kodum.
}

func respondWithMessage(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
} // payload olarak message: ya da error: veriyon.

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
