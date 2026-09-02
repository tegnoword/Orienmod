package com.example.orinmodapp.data.api


import com.apollographql.apollo.ApolloClient
import com.apollographql.apollo.network.http.LoggingInterceptor

object ApolloClientProvider {
    private const val BASE_URL = "http://localhost:8080/query?"
    val apolloClient: ApolloClient by lazy {
        ApolloClient.Builder()
            .serverUrl(BASE_URL)
            .addHttpInterceptor(
                LoggingInterceptor(
                    level = LoggingInterceptor.Level.BODY,
                    log = { message -> println("Apollo: $message") }
                )
            )
            .build()
    }
}