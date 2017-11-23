package chameleon

import (
    "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func mobileToTeam(c context.Context, mobile string) *datastore.Key
{
    q := datastore.NewQuery("TeamMember").
        Filter("number =", mobile).
        KeysOnly().
        Limit(1)
    for t := q.Run(ctx); ; {
        key, err := t.Next(nil)
        if err != nil {
            break
        }
        return key.Parent()
    }
    return nil
}
