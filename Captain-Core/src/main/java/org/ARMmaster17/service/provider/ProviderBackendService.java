package org.ARMmaster17.service.provider;

public interface ProviderBackendService {

    /**
     * Creates a machine (VM/container/baremetal) on the underlying provider infrastructure.
     * @return Provider's unique identifier for this machine (if the provisioning succeeds).
     */
    String CreateMachine(String hostname, int cores, int memory, String ip, int storageSize) throws Exception;

    void DestroyMachine(String providerID) throws Exception;

    void Ping() throws Exception;
}
