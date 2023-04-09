package ru.itis.products;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import ru.itis.products.service.ProductServiceClient;

import java.util.Scanner;

@SpringBootApplication
public class JavaClientGrpcApplication implements CommandLineRunner {

	@Autowired
	private ProductServiceClient productServiceClient;

	public static void main(String[] args) {
		SpringApplication.run(JavaClientGrpcApplication.class, args);
	}

	@Override
	public void run(String... args) throws Exception {
		Scanner scanner = new Scanner(System.in);
		String id = scanner.nextLine();

		System.out.println(productServiceClient.getNameOfProduct(id));
	}
}
