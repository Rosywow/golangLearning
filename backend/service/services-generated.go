// Package services generated by service-enroll-generate.go;
// !DO NOT EDIT THIS FILE!
package service

import (
	"w2w.io/serve/authmgmt" //auth-mgmt ,  auth, , XUnion@GMail.com
	"w2w.io/serve/document" //document ,  document, 13580452503, KManager@GMail.com
	"w2w.io/serve/logview"  //log-view ,  log-view, 18928776452, XUnion@GMail.com
	"w2w.io/serve/message"  //message-mgr ,  tom sawyer, 13580452503, KManager@GMail.com
	"w2w.io/serve/nservice" //newsvc ,  newservice, , aaa@GMail.com
	"w2w.io/serve/static"   //static ,  static, 18928776452, XUnion@GMail.com
	"w2w.io/serve/user"     //user-mgmt ,  user, 18928776452, XUnion@GMail.com
)

// Enroll will be called from serve cmd
func Enroll() {

	authmgmt.Enroll(`{"name":"auth","email":"XUnion@GMail.com"}`)
	document.Enroll(`{"name":"document","tel":"13580452503","email":"KManager@GMail.com"}`)
	logview.Enroll(`{"name":"log-view","tel":"18928776452","email":"XUnion@GMail.com"}`)
	message.Enroll(`{"name":"tom sawyer","tel":"13580452503", "email":"KManager@GMail.com"}`)
	nservice.Enroll(`{"name":"newservice","email":"aaa@GMail.com"}`)
	static.Enroll(`{"name":"static","tel":"18928776452","email":"XUnion@GMail.com"}`)
	user.Enroll(`{"name":"user","tel":"18928776452","email":"XUnion@GMail.com"}`)
}
