# Hydra host container #

## Build the container ##

    make

## Startup ##

Start `hydra-postgres` first.

Then start the container:

    docker run -d -p 9090:9090 --net hydra --name hydra-host hydra-host

If you use `docker-machine`, redirect port 9090.
The server should now be started. You can check by opening the address `http://localhost:9090/alive`

## Access the Hydra command line ##

Once the container is started, you can call the CLI interface using `docker exec`:

    docker docker exec -it hydra-host /hydra-host <command>

Suggestion: create an alias to be able to simply use the `hydra-host` command, and copy-paste the examples:

    alias hydra-host='docker exec -it hydra-host /hydra-host'

## Examples ##

### Create the superuser ###

hydra-host account create [command options] <username>

OPTIONS:
   --password           the user's password
   --as-superuser       grant superuser privileges to the user
so for superuser :
hydra-host account create --as-superuser  --password secret admin

#### to do the authent of this user you need a generic client app ####
    hydra-host client create -i generic-client -s secretgeneric -r http://localhost:8080/authenticate/callback

#### Authenticate the superuser with this generic client ####
    curl -kd 'grant_type=password&username=admin&password=secret' -u generic-client:secretgeneric http://localhost:9090/oauth2/token
    You should get his token

### Create a super user client (application) ###

Create a super-user client with id `app` and password `secret` :

    hydra-host client create -i app -s secret -r http://localhost:8080/authenticate/callback --as-superuser

### Get an access token for a client ###

    curl -k -X POST --user app:secret "http://localhost:9090/oauth2/token?grant_type=client_credentials"

Copy the `access_token` field for later use.

### Create a user ###

Create a user `foo@bar.com` with password `secret`:

    curl -ksSd '{"username": "foo@bar.com", "password": "secret"}' -H "Authorization: Bearer <client token>" http://localhost:9090/accounts

Replace `<client token>` with the token from earlier.

### Create a policy ###

Create a policy named "Resource 1 get" to `get` the resource `test:resource1` with all users:

    curl -kd '{"description": "Resource 1 get", "subjects": ["<.*>"], "resources": ["test:resource1"], "effect": "allow", "permissions": ["get"]}' -u app:secret  http://localhost:9090/policies
