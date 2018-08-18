# gmail-receiver-golang
It is a gmail receiver tool using google APIs, written in Golang. 

For accessing your Gmail, a google API credential must be provided, which can be configured from the following link for a specific google account that has a gmail affiliated. https://console.developers.google.com/flows/enableapi?apiid=gmail

	1. Activate the google api as from link above for your own gmail. Rename the finalized credential as: client_secret.json, to put into the same folder.
	2. 	# go get -u google.golang.org/api/gmail/v1 
		# go get -u golang.org/x/oauth2/...
	2. To include a 3rd party crypto lib: 
		# go get github.com/Luzifer/go-openssl  (You may not need to “go get” these libs since I have packed all dependencies into the local vendor folder.)
	3. 	# go build receive.go
	4. 	# ./receive -credential client_secret.json -from 2018/01/01 -to 2018/04/10 -enc -enckey password
		(Either plaintext email or encrypted email can be saved locally by using “-enc” flag. Can check more flag explanation by typing: ./receive --help in console)
    
The work is based on the official tutorial from Google, and added with more, like screening an encryption, options. 
Do give your own client_secret.json for testing.
