## Mixpanel Client in Go

Refer to:
[http://godoc.org/github.com/nitrous-io/go-mixpanel](http://godoc.org/github.com/nitrous-io/go-mixpanel)
for documentation.

Example code here:

```go
package main

import "github.com/nitrous-io/go-mixpanel"

func main() {
	mc := mixpanel.NewMixpanelClient("your_mixpanel_token")
  err := mc.CreateProfile("deadbeef", map[string]interface{}{"$first_name": "Mclovin"})
	if err != nil {
		panic(err)
	}

	err = mc.Alias("deadbeef", "1")
	if err != nil {
		panic(err)
	}

	err = mc.IncrementPropertiesOnProfile("1", map[string]int{"hosts_created": 1})
	if err != nil {
		panic(err)
  }

  err = mc.Track("User Signed Up", map[string]interface{}{"$distinct_id":"1"})
  if err != nil {
    panic(err)
  }
}
```
