package handler

type NodeType int

const (
	SOURCE NodeType = iota + 1
	TARGET
	UNKNOWN
)

// func NodeSearchQuery(c *gin.Context) {
// 	lng, err := strconv.ParseFloat(c.Query("lng"))
// 	checkErr(err)
// 	lat, err := strconv.ParseFloat(c.Query("lat"))
// 	checkErr(err)
// 	query, err := c.Query("query")
// 	checkErr(err)

// 	fmt.Printf("calling NodeSearchQuery lng=%s lat=%s query=%s", lng, lat.query)

// 	c.JSON(http.StatusOK, edges)
// }

// @GET
// @Path("/find")
// @Produces(MediaType.APPLICATION_JSON)
// public Response getSearchWithQueryString(@QueryParam(value = "lng")
// String lng, @QueryParam(value = "lat")
// String lat, @QueryParam(value = "query")
// String queryString) throws WebGISApplicationException
// {
// 	Integer result = this.nodeService.findNode(lng, lat, NodeService.TYPE.UNKNOWN);
// 	if(result == null)
// 	{
// 		return Response.status(Status.NOT_FOUND).build();
// 	}
// 	return Response.ok().entity(result).build();

// }

// @GET
// @Path("/find/source")
// @Produces(MediaType.APPLICATION_JSON)
// public Response getSearchSource(@QueryParam(value = "lng")
// String lng, @QueryParam(value = "lat")
// String lat, @QueryParam(value = "query")
// String queryString) throws WebGISApplicationException
// {
// 	Integer result = this.nodeService.findNode(lng, lat, NodeService.TYPE.SOURCE);
// 	if(result == null)
// 	{
// 		return Response.status(Status.NOT_FOUND).build();
// 	}
// 	return Response.ok().entity(result).build();

// }

// @GET
// @Path("/find/target")
// @Produces(MediaType.APPLICATION_JSON)
// public Response getSearchTarget(@QueryParam(value = "lng")
// String lng, @QueryParam(value = "lat")
// String lat, @QueryParam(value = "query")
// String queryString) throws WebGISApplicationException
// {
// 	Integer result = this.nodeService.findNode(lng, lat, NodeService.TYPE.TARGET);
// 	if(result == null)
// 	{
// 		return Response.status(Status.NOT_FOUND).build();
// 	}
// 	return Response.ok().entity(result).build();

// }

// -------

// get lat=, lon=

// switch(type)
// 			{
// 				case SOURCE:
// 					filter = "where eout > 0";
// 					break;
// 				case TARGET:
// 					filter = "where ein > 0";
// 					break;
// 				default:
// 					filter = "";
// 					break;
// 			}

// 			String sql = "select id from public.roads_vertices_pgr " + filter + " order by geom <-> ST_GeomFromText( ? , 4326) limit 1";

// 			String lng, String lat

// 	Query 1		Point(lon,lat)
