// Copyright (C) 2021-2022 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package apiv1

import (
	"git.lcomrade.su/root/lenpaste/internal/storage"
	"encoding/json"
	"net/http"
)

// GET /api/v1/get?id=""
func (data Data) GetHand(rw http.ResponseWriter, req *http.Request) {
	// Get paste id
	req.ParseForm()
	
	id := req.PostForm.Get("id")

	// Check paste id
	if id == "" {
		data.writeError(rw, req, storage.ErrNotFoundID)
	}

	// Get paste
	paste, err := data.DB.PasteGet(id)
	if err != nil {
		data.writeError(rw, req, err)
		return
	}

	// Return response
	rw.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(rw).Encode(paste)
	if err != nil {
		data.Log.HttpError(req, err)
		return
	}
}
