package org.ARMmaster17.resource;

import io.quarkus.hibernate.orm.panache.PanacheEntityBase;
import io.quarkus.hibernate.orm.panache.PanacheQuery;
import org.ARMmaster17.entity.Plane;

import javax.ws.rs.*;
import javax.ws.rs.core.Response;

@Path("/planes")
public class PlaneResource {
    @GET
    @Produces("application/json")
    public Response list() {
        PanacheQuery<PanacheEntityBase> planes = Plane.findAll();
        return Response.ok(planes.list()).build();
    }
}
