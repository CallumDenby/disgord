package disgord

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andersfylling/disgord/httd"
	"github.com/andersfylling/snowflake/v3"
)

type ErrRest = httd.ErrREST

// URLParameters converts a struct of values to a valid URL query string
type URLParameters interface {
	GetQueryString() string
}

func unmarshal(data []byte, v interface{}) error {
	return httd.Unmarshal(data, v)
}

func marshal(v interface{}) ([]byte, error) {
	return httd.Marshal(v)
}

// AvatarParamHolder is used when handling avatar related REST structs.
// since a Avatar can be reset by using nil, it causes some extra issues as omit empty cannot be used
// to get around this, the struct requires an internal state and must also handle custom marshalling
type AvatarParamHolder interface {
	json.Marshaler
	Empty() bool
	SetAvatar(avatar string)
	UseDefaultAvatar()
}

func newRESTBuilder(cache *Cache, client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) *RESTBuilder {
	builder := &RESTBuilder{}
	builder.setup(cache, client, config, middleware)

	return builder
}

type paramHolder map[string]interface{}

func (p paramHolder) GetQueryString() string {
	if len(p) == 0 {
		return ""
	}

	var params string
	seperator := "?"
	for k, v := range p {
		var str string
		var uintHolder uint64
		switch t := v.(type) {
		case snowflake.ID:
			str = t.String()
		case int8:
			uintHolder = uint64(t)
		case int16:
			uintHolder = uint64(t)
		case int32:
			uintHolder = uint64(t)
		case int64:
			uintHolder = uint64(t)
		case int:
			uintHolder = uint64(t)
		case uint8:
			uintHolder = uint64(t)
		case uint16:
			uintHolder = uint64(t)
		case uint32:
			uintHolder = uint64(t)
		case uint:
			uintHolder = uint64(t)
		case uint64:
			uintHolder = t
		case bool:
			if t {
				str = "true"
			} else {
				str = "false"
			}
		case string:
			str = t
		}
		// TODO float
		if str == "" {
			str = strconv.FormatUint(uintHolder, 10)
		}

		params += seperator + k + "=" + str
		seperator = "&"
	}

	return params
}

var _ URLParameters = (*paramHolder)(nil)

type fRESTRequestMiddleware func(resp *http.Response, body []byte, err error) error
type fRESTCacheMiddleware func(resp *http.Response, v interface{}, err error) error
type fRESTItemFactory func() interface{}

//go:generate go run generate/restbuilders/main.go

type RESTBuilder struct {
	middleware fRESTRequestMiddleware
	config     *httd.Request
	client     httd.Requester

	itemFactory fRESTItemFactory

	cache           *Cache
	cacheRegistry   cacheRegistry
	cacheMiddleware fRESTCacheMiddleware
	cacheItemID     snowflake.ID

	body              map[string]interface{}
	urlParams         paramHolder
	ignoreCache       bool
	cancelOnRatelimit bool
}

func (b *RESTBuilder) setup(cache *Cache, client httd.Requester, config *httd.Request, middleware fRESTRequestMiddleware) {
	b.body = make(map[string]interface{})
	b.urlParams = make(map[string]interface{})
	b.cache = cache
	b.client = client
	b.config = config
	b.middleware = middleware

	if b.config == nil {
		b.config = &httd.Request{
			Method: http.MethodGet,
		}
	}
}

func (b *RESTBuilder) cacheLink(registry cacheRegistry, middleware fRESTCacheMiddleware) {
	b.cacheRegistry = registry
	b.cacheMiddleware = middleware
}

func (b *RESTBuilder) prepare() {
	// update the config
	if b.config.ContentType != "" {
		b.config.Body = b.body
	}
	b.config.Endpoint += b.urlParams.GetQueryString()

	if b.cache == nil {
		b.IgnoreCache()
	}
}

// execute ... v must be a nil pointer.
func (b *RESTBuilder) execute() (v interface{}, err error) {
	if !b.ignoreCache && b.config.Method == http.MethodGet && !b.cacheItemID.Empty() {
		// cacheLink lookup. return on cacheLink hit
		v, err = b.cache.Get(b.cacheRegistry, b.cacheItemID)
		if v != nil && err == nil {
			return v, nil
		}
		// otherwise we perform the request
	}

	b.prepare()

	var resp *http.Response
	var body []byte
	resp, body, err = b.client.Request(b.config)
	if err != nil {
		return nil, err
	}

	if b.middleware != nil {
		if err = b.middleware(resp, body, err); err != nil {
			return nil, err
		}
	}

	if len(body) > 1 {
		v = b.itemFactory()
		if err = httd.Unmarshal(body, v); err != nil {
			return nil, err
		}

		if b.cacheRegistry == NoCacheSpecified {
			return v, err
		}

		if b.cacheMiddleware != nil {
			b.cacheMiddleware(resp, v, err)
		}

		b.cache.Update(b.cacheRegistry, v)
	}
	return v, nil
}

type restReqBuilderAsync struct {
	Data interface{}
	Err  error
	// FromCache bool // TODO
}

func (b *RESTBuilder) async() <-chan *restReqBuilderAsync {
	A := make(chan *restReqBuilderAsync)
	go func() {
		resp := &restReqBuilderAsync{}
		resp.Data, resp.Err = b.execute()

		A <- resp
	}()

	return A
}

func (b *RESTBuilder) param(name string, v interface{}) *RESTBuilder {
	if b.config.Method == http.MethodGet {
		// RFC says you can not send a body in a GET request
		b.queryParam(name, v)
	} else {
		b.body[name] = v
	}
	return b
}

func (b *RESTBuilder) queryParam(name string, v interface{}) *RESTBuilder {
	b.urlParams[name] = v
	return b
}

func (b *RESTBuilder) IgnoreCache() *RESTBuilder {
	b.ignoreCache = true
	return b
}

func (b *RESTBuilder) CancelOnRatelimit() *RESTBuilder {
	b.cancelOnRatelimit = true
	return b
}

// GetGateway [REST] Returns an object with a single valid WSS URL, which the client can use for Connecting.
// Clients should cacheLink this value and only call this endpoint to retrieve a new URL if they are unable to
// properly establish a connection using the cached version of the URL.
//  Method                  GET
//  Endpoint                /gateway
//  Rate limiter            /gateway
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway
//  Reviewed                2018-10-12
//  Comment                 This endpoint does not require authentication.
func GetGateway(client httd.Getter) (gateway *Gateway, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: "/gateway",
		Endpoint:    "/gateway",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}

// GetGatewayBot [REST] Returns an object based on the information in Get Gateway, plus additional metadata
// that can help during the operation of large or sharded bots. Unlike the Get Gateway, this route should not
// be cached for extended periods of time as the value is not guaranteed to be the same per-call, and
// changes as the bot joins/leaves guilds.
//  Method                  GET
//  Endpoint                /gateway/bot
//  Rate limiter            /gateway/bot
//  Discord documentation   https://discordapp.com/developers/docs/topics/gateway#get-gateway-bot
//  Reviewed                2018-10-12
//  Comment                 This endpoint requires authentication using a valid bot token.
func GetGatewayBot(client httd.Getter) (gateway *GatewayBot, err error) {
	var body []byte
	_, body, err = client.Get(&httd.Request{
		Ratelimiter: "/gateway/bot",
		Endpoint:    "/gateway/bot",
	})
	if err != nil {
		return
	}

	err = unmarshal(body, &gateway)
	return
}
