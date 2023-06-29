# The OpenStack Barbican SDK for Golang

## Usage

```go
import (
	"context"
	"fmt"

	barbican "github.com/artashesbalabekyan/barbican-sdk-go"
	"github.com/artashesbalabekyan/barbican-sdk-go/xhttp"
)

func main() {
	ctx := context.Background()

	config := &xhttp.Config{
		Endpoint: "https://<endpoint>",
		Login: xhttp.Credentials{
			ProjectDomain:  "default",
			ProjectName:    "<project_name>",
			AuthUrl:        "https://<auth_url>",
			Username:       "<userName>",
			Password:       "<password>",
			UserDomainName: "Default",
		},
	}

	client, err := barbican.NewConnection(ctx, config)
	if err != nil {
		panic(err)
	}
}

```

```go
err := client.Create(ctx, "my-key", []byte("my-value"))
if err != nil {
    panic(err)
}

```

```go
iterator, err := client.ListSecrets(ctx)
if err != nil {
    panic(err)
}

defer iterator.Close()

for {
    name, ok := iterator.Next()
    if !ok {
        break
    }
    fmt.Println(name)
}
```