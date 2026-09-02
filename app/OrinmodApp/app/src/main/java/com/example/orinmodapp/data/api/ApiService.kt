package com.example.orinmodapp.data.api

import com.apollographql.apollo.ApolloClient
import com.orienmod.app.graphql.CoursesQuery

class GraphQLService {

    private val apolloClient: ApolloClient = ApolloClientProvider.apolloClient

    suspend fun getCourses(email: String): Result<List<CoursesQuery.Course>> {
        return try {
            val response = apolloClient.query(CoursesQuery(email = email)).execute()
            if (response.hasErrors()) {
                val errors = response.errors?.joinToString { it.message } ?: "Error desconocido"
                Result.failure(Exception(errors))
            } else {
                Result.success(response.data?.courses ?: emptyList())
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
