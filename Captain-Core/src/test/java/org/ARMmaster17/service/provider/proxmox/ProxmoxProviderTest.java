package org.ARMmaster17.service.provider.proxmox;

import io.quarkus.test.junit.QuarkusTest;
import org.ARMmaster17.service.provider.proxmox.ProxmoxProvider;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import javax.inject.Inject;

@QuarkusTest
public class ProxmoxProviderTest {
    @Inject
    ProxmoxProvider proxmoxProvider;

    @Test
    public void testCreateMachine() throws Exception {
        String hostname = "test";
        int cores = 1;
        int memory = 1024;
        int storage = 10;
        String ip = "10.9.0.1/8";
        String result = proxmoxProvider.CreateMachine(hostname, cores, memory, ip, storage);
        Assertions.assertNotNull(result);
        Assertions.assertEquals("106", result);
    }

    @Test
    public void testDestroyMachine() throws Exception {
        proxmoxProvider.DestroyMachine("106");
    }

    @Test
    public void testGetNextID() throws Exception {
        Assertions.assertEquals(120, proxmoxProvider.getNextVMID());
    }
}
