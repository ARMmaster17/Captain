package org.ARMmaster17.service.provider.proxmox;

import io.quarkus.test.junit.QuarkusTest;
import it.corsinvest.proxmoxve.api.PveClient;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import javax.inject.Inject;

@QuarkusTest
public class ProxmoxVEClientConfigTest {
    @Inject
    PveClient pveClient;

    @Test
    public void testClientIsConnected() {
        Assertions.assertEquals(3, pveClient.getNodes().index().getResponse().getJSONArray("data").length());
    }
}
