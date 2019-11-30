const { parse } = require('graphql');

const FILE_EXTENSION = 'go';

const routeTemplate = (mutationSdl, typesSdl) => {
  const mutationAst = parse(mutationSdl);
  const actionName = mutationAst.definitions[0].fields[0].name.value;
  return `
// handler.HandleFunc("/${actionName}", ${actionName}Handler)
`;
};


const handlerTemplate = (
  mutationSdl,
  typesSdl
) => {
  const types = parse(`${typesSdl}\n${mutationSdl}`);

  const actionDefinition = types.definitions.find(t => t.name.value === 'Mutation').fields[0];
  const actionName = actionDefinition.name.value;

  return `
package main

import (
	"encoding/json"
	"net/http"
)

var ${actionName}Handler = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// do your business logic here

	var res responseBody
	res.Data = map[string]interface{}{
		"id": 1,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
`;
};
