package org.ARMmaster17.service;

import org.ARMmaster17.entity.Plane;
import org.ARMmaster17.service.provider.ProviderBackendService;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

@ApplicationScoped
public class PlaneService {

    @Inject
    ProviderBackendService providerBackendService;

    public Plane CreatePlane(String hostname, int instanceNumber, int cores, int memory, int storageSize, String ip) throws Exception {
        String numberedHostname = String.format("%s%d", hostname, instanceNumber);
        return CreatePlane(numberedHostname, cores, memory, storageSize, ip);
    }

    public Plane CreatePlane(String hostname, int cores, int memory, int storageSize, String ip) throws Exception {
        Plane plane = new Plane();
        plane.Hostname = hostname;
        plane.Cores = cores;
        plane.Memory = memory;
        plane.IP = ip;
        plane.StorageSize = storageSize;
        plane.ProviderIdentifier = providerBackendService.CreateMachine(plane.Hostname, plane.Cores, plane.Memory, plane.IP, plane.StorageSize);
        plane.persist();
        return plane;
    }

    public void DestroyPlane(Plane plane) throws Exception {
        providerBackendService.DestroyMachine(plane.ProviderIdentifier);
        plane.delete();
    }
}
