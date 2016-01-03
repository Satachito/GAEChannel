package main

import (
	"net/http"
	"html/template"

	"appengine"
	"appengine/channel"
	"appengine/datastore"
)

func
init() {
    http.HandleFunc( "/", handle )
	http.HandleFunc( "/_ah/channel/connected/", connected )
	http.HandleFunc( "/_ah/channel/disconnected/", disconnected )
}

func
ChannelKey( ctx appengine.Context, id string ) *datastore.Key {
	return datastore.NewKey( ctx, "Channels", id, 0, nil )
}

type
NULL struct {
}

func
connected( w http.ResponseWriter, r *http.Request ) {
	ctx := appengine.NewContext( r )

	if _, err := datastore.Put(
		ctx,
		ChannelKey( ctx, r.FormValue( "from" ) ),
		&NULL{},
	); err != nil { panic( err ) }
}

func
disconnected( w http.ResponseWriter, r *http.Request ) {
	ctx := appengine.NewContext( r )

	if err := datastore.Delete(
		ctx,
		ChannelKey( ctx, r.FormValue( "from" ) ),
	); err != nil { panic( err ) }
}

func
Execute( tag string, w http.ResponseWriter, p interface{} ) {
	if err := template.Must( template.ParseFiles( tag + ".html" ) ).Execute(
		w,
		p,
	); err != nil { panic( err ) }
}

func
handle( w http.ResponseWriter, r *http.Request ) {
	switch r.Method {
	case "GET":
		Execute( "get", w, nil )
	case "POST":
		wID := r.FormValue( "ID" )
		wToken, err := channel.Create( appengine.NewContext( r ), wID )
		if err != nil { 
			Execute( "get", w, map[ string ]string { "Error": err.Error() } )
		} else {
			Execute(
				"post",
				w,
				map[ string ]string {
					"ID"		: wID,
					"Token"		: wToken,
				},
			)
		}
	case "PUT":
		ctx := appengine.NewContext( r )

		keys, err := datastore.NewQuery( "Channels" ).KeysOnly().GetAll( ctx, nil )
		if err != nil { panic( err ) }

		for _, key := range keys {
			if err := channel.SendJSON(
				ctx,
				key.StringID(),
				map[ string ]interface{} {
					"message"	: r.FormValue( "Message" ),
					"ID"		: r.FormValue( "ID" ),
				},
			); err != nil { panic( err ) }
		}
	}
}

