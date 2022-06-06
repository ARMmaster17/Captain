package org.ARMmaster17.service.provider;

import org.ARMmaster17.service.provider.proxmox.ProxmoxProvider;

import javax.enterprise.context.ApplicationScoped;
import javax.ws.rs.Produces;

@ApplicationScoped
public class ProviderService {
    @Produces
    public ProviderBackendService providerBackendService() {
        return new ProxmoxProvider();
    }
}
