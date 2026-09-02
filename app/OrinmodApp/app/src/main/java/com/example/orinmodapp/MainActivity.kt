package com.example.orinmodapp

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.example.orinmodapp.data.api.GraphQLService
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()
        setContent {
            TestConnectionScreen()
        }
    }
}

@Composable
fun TestConnectionScreen(modifier: Modifier = Modifier) {
    var statusText by remember { mutableStateOf("Presiona el botón para probar la conexión") }
    var isLoading by remember { mutableStateOf(false) }
    var coursesList by remember { mutableStateOf<List<String>>(emptyList()) }

    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Button(
            onClick = {
                isLoading = true
                statusText = "Probando conexión..."
                coursesList = emptyList()

                CoroutineScope(Dispatchers.Main).launch {
                    try {
                        val service = GraphQLService()
                        val result = withContext(Dispatchers.IO) {
                            service.getCourses("test@example.com") // Cambia por un email real
                        }

                        result.onSuccess { courses ->
                            statusText = "✅ ${courses.size} cursos encontrados"
                            coursesList = courses.map { "${it.name} (${it.id})" }
                        }.onFailure { error ->
                            statusText = "❌ Error: ${error.message}"
                        }
                    } catch (e: Exception) {
                        statusText = "💥 Error: ${e.message}"
                    } finally {
                        isLoading = false
                    }
                }
            },
            enabled = !isLoading
        ) {
            Text(if (isLoading) "Cargando..." else "Probar Conexión")
        }
    }

    // Mostrar el estado en la parte inferior
    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.BottomCenter
    ) {
        Text(
            text = statusText,
            modifier = Modifier.padding(bottom = 50.dp)
        )
    }

    // Mostrar lista de cursos si hay
    if (coursesList.isNotEmpty()) {
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Column (
                horizontalAlignment = Alignment.CenterHorizontally,
                modifier = Modifier.padding(top = 100.dp)
            ) {
                coursesList.forEach { course ->
                    Text(text = "📚 $course")
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun TestConnectionPreview() {
    TestConnectionScreen()
}