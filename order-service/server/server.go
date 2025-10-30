package server

import (
    "html/template"
    "log"
    "net/http"
    "order-service/cache"
  

    "github.com/gorilla/mux"
)

var (
    orderTmpl = template.Must(template.ParseFiles("templates/order.html"))
    listTmpl  = template.Must(template.ParseFiles("templates/list.html"))
)

func StartServer(cache *cache.OrderCache) {
    r := mux.NewRouter()


    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        orders := cache.List()
        listTmpl.Execute(w, orders)
    })

   
    r.HandleFunc("/order/{id}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        order, ok := cache.Get(id)
        if !ok {
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte("<h1>Order Not Found</h1><p><a href='/'>← На главную</a></p>"))
            return
        }

        orderTmpl.Execute(w, order)
    })

    log.Println("HTTP server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))

}
