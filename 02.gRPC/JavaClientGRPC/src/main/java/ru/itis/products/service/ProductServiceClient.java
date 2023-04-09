package ru.itis.products.service;

import net.devh.boot.grpc.client.inject.GrpcClient;
import org.springframework.stereotype.Service;
import ru.itis.products.pb.ProductRequest;
import ru.itis.products.pb.ProductServiceGrpc;

@Service
public class ProductServiceClient {

    @GrpcClient("product-service")
    private ProductServiceGrpc.ProductServiceBlockingStub service;

    public String getNameOfProduct(String id) {
        return service.getProduct(ProductRequest.newBuilder()
                .setId(id)
                .build()).getName();
    }
}
