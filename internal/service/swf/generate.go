//go:generate go run ../../generate/tags/main.go -AWSSDKVersion=2 -ListTags -ServiceTagsSlice -TagType=ResourceTag -UpdateTags
//go:generate go run ../../generate/servicepackage/main.go
// ONLY generate directives and package declaration! Do not add anything else to this file.

package swf
