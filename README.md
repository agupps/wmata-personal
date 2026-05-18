Front end is deployed via vercel: https://wmata-personal.vercel.app/

Back end is deployed via Render: https://wmata-personal.onrender.com

Endpoints for backend are /metroStops and /busStops

/busStops
Query Params: ?lines=D72&stop=1001694 (the one closes to my house)

/metroStops
Query Params: ?lines=SV&stop=C03 (farragut west)

To run the backend: 

```
go run cmd/main.go
```

To run the front end locally:

```
npm run dev
```

To build 
```
npm run build
```

