# Examples

This directory contains examples for how to use zpages with the various go web
frameworks that exist.

* basic/ - This example shows how to use zpages with no framework.
* chi/ - This eample shows how to use zpages with the [chi framework](https://github.com/go-chi/chi).
* gorilla/ - This example shows how to use zpages with the [gorilla mux framework](https://github.com/gorilla/mux).
* gin/ - This example shows how to use zpages with the [gin framework](https://github.com/gin-gonic/gin). Unfortunately, the Requestz doesn't work with gin right now since the http.ResponseWriter status code capturing doesn't carry over across the gin.Context boundary.
