package disgord

// Warning: This file has been automatically generated by generate/restbuilders/main.go
// DO NOT EDIT! This file is overwritten at "go generate"
// This file holds all the basic RESTBuilder methods a builder is expected to.

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *guildAuditLogsBuilder) IgnoreCache() *guildAuditLogsBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *guildAuditLogsBuilder) CancelOnRatelimit() *guildAuditLogsBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *guildAuditLogsBuilder) URLParam(name string, v interface{}) *guildAuditLogsBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *guildAuditLogsBuilder) Set(name string, v interface{}) *guildAuditLogsBuilder {
	b.r.body[name] = v
	return b
}

func (b *guildAuditLogsBuilder) SetUserID(userID Snowflake) *guildAuditLogsBuilder {
	b.r.param("user_id", userID)
	return b
}

func (b *guildAuditLogsBuilder) SetActionType(actionType uint) *guildAuditLogsBuilder {
	b.r.param("action_type", actionType)
	return b
}

func (b *guildAuditLogsBuilder) SetBefore(before Snowflake) *guildAuditLogsBuilder {
	b.r.param("before", before)
	return b
}

func (b *guildAuditLogsBuilder) SetLimit(limit int) *guildAuditLogsBuilder {
	b.r.param("limit", limit)
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *listGuildEmojisBuilder) IgnoreCache() *listGuildEmojisBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *listGuildEmojisBuilder) CancelOnRatelimit() *listGuildEmojisBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *listGuildEmojisBuilder) URLParam(name string, v interface{}) *listGuildEmojisBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *listGuildEmojisBuilder) Set(name string, v interface{}) *listGuildEmojisBuilder {
	b.r.body[name] = v
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *deleteInviteBuilder) IgnoreCache() *deleteInviteBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *deleteInviteBuilder) CancelOnRatelimit() *deleteInviteBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *deleteInviteBuilder) URLParam(name string, v interface{}) *deleteInviteBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *deleteInviteBuilder) Set(name string, v interface{}) *deleteInviteBuilder {
	b.r.body[name] = v
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *getInviteBuilder) IgnoreCache() *getInviteBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *getInviteBuilder) CancelOnRatelimit() *getInviteBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *getInviteBuilder) URLParam(name string, v interface{}) *getInviteBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *getInviteBuilder) Set(name string, v interface{}) *getInviteBuilder {
	b.r.body[name] = v
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *modifyGuildRoleBuilder) IgnoreCache() *modifyGuildRoleBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *modifyGuildRoleBuilder) CancelOnRatelimit() *modifyGuildRoleBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *modifyGuildRoleBuilder) URLParam(name string, v interface{}) *modifyGuildRoleBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *modifyGuildRoleBuilder) Set(name string, v interface{}) *modifyGuildRoleBuilder {
	b.r.body[name] = v
	return b
}

func (b *modifyGuildRoleBuilder) SetName(name string) *modifyGuildRoleBuilder {
	b.r.param("name", name)
	return b
}

func (b *modifyGuildRoleBuilder) SetPermissions(permissions uint64) *modifyGuildRoleBuilder {
	b.r.param("permissions", permissions)
	return b
}

func (b *modifyGuildRoleBuilder) SetColor(color uint) *modifyGuildRoleBuilder {
	b.r.param("color", color)
	return b
}

func (b *modifyGuildRoleBuilder) SetHoist(hoist bool) *modifyGuildRoleBuilder {
	b.r.param("hoist", hoist)
	return b
}

func (b *modifyGuildRoleBuilder) SetMentionable(mentionable bool) *modifyGuildRoleBuilder {
	b.r.param("mentionable", mentionable)
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *getUserBuilder) IgnoreCache() *getUserBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *getUserBuilder) CancelOnRatelimit() *getUserBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *getUserBuilder) URLParam(name string, v interface{}) *getUserBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *getUserBuilder) Set(name string, v interface{}) *getUserBuilder {
	b.r.body[name] = v
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *putUserBuilder) IgnoreCache() *putUserBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *putUserBuilder) CancelOnRatelimit() *putUserBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *putUserBuilder) URLParam(name string, v interface{}) *putUserBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *putUserBuilder) Set(name string, v interface{}) *putUserBuilder {
	b.r.body[name] = v
	return b
}

func (b *putUserBuilder) SetUsername(username string) *putUserBuilder {
	b.r.param("username", username)
	return b
}

func (b *putUserBuilder) SetAvatar(avatar string) *putUserBuilder {
	b.r.param("avatar", avatar)
	return b
}

// IgnoreCache will not fetch the data from the cache if available, and always execute a
// a REST request. However, the response will always update the cache to keep it synced.
func (b *listVoiceRegionsBuilder) IgnoreCache() *listVoiceRegionsBuilder {
	b.r.IgnoreCache()
	return b
}

// CancelOnRatelimit will disable waiting if the request is rate limited by Discord.
func (b *listVoiceRegionsBuilder) CancelOnRatelimit() *listVoiceRegionsBuilder {
	b.r.CancelOnRatelimit()
	return b
}

// URLParam adds or updates an existing URL parameter.
// eg. URLParam("age", 34) will cause the URL `/test` to become `/test?age=34`
func (b *listVoiceRegionsBuilder) URLParam(name string, v interface{}) *listVoiceRegionsBuilder {
	b.r.queryParam(name, v)
	return b
}

// Set adds or updates an existing a body parameter
// eg. Set("age", 34) will cause the body `{}` to become `{"age":34}`
func (b *listVoiceRegionsBuilder) Set(name string, v interface{}) *listVoiceRegionsBuilder {
	b.r.body[name] = v
	return b
}