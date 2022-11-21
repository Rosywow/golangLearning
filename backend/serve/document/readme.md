# static files serve

Set Ep.Path to deployed url path, for example
   /images/

Set Ep.DocRoot to the physical disk path, for example
    /tmp/images/

## 1 common files
  set Ep.PageRoute to false, so, if the file missing then browser will get 404


## 2 javascript SPA
Set Ep.PageRoute to true, so, if the request path nonexistence, serve will redirect to Ep.DocRoot/**/index.html

## 3 customize serve api
If Ep.Fn unset then cmn.AddService will use the default cmn.WebFS func, if set then use the set value as serve api.


<style>
h4{
    color:#0aceeb;
}

h3{
    color:#007bac;
}

h2{
    color:#01a3b0;
}

h1{
    color:#5d8aa8;
}
</style>