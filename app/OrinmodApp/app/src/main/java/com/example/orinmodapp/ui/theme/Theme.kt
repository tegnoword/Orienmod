package com.example.orinmodapp.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat

// ============================================
// IMPORTAR TYPOGRAPHY STYLES
// ============================================
// ✅ Estos vienen de Typography.kt
import com.example.orinmodapp.ui.theme.HeadlineXL
import com.example.orinmodapp.ui.theme.HeadlineLG
import com.example.orinmodapp.ui.theme.HeadlineLGMobile
import com.example.orinmodapp.ui.theme.TitleMD
import com.example.orinmodapp.ui.theme.BodyLG
import com.example.orinmodapp.ui.theme.BodyMD
import com.example.orinmodapp.ui.theme.LabelMD
import com.example.orinmodapp.ui.theme.LabelSM

// ============================================
// LIGHT COLOR SCHEME
// ============================================
val LightColorScheme = lightColorScheme(
    primary = Color(0xFF6750A4),        // ⚠️ Reemplaza con tus colores
    onPrimary = Color(0xFFFFFFFF),
    primaryContainer = Color(0xFFEADDFF),
    onPrimaryContainer = Color(0xFF21005D),
    secondary = Color(0xFF625B71),
    onSecondary = Color(0xFFFFFFFF),
    secondaryContainer = Color(0xFFE8DEF8),
    onSecondaryContainer = Color(0xFF1D192B),
    tertiary = Color(0xFF7D5260),
    onTertiary = Color(0xFFFFFFFF),
    tertiaryContainer = Color(0xFFFFD8E4),
    onTertiaryContainer = Color(0xFF31111D),
    error = Color(0xFFBA1A1A),
    onError = Color(0xFFFFFFFF),
    errorContainer = Color(0xFFFFDAD6),
    onErrorContainer = Color(0xFF410002),
    background = Color(0xFFFFFBFE),
    onBackground = Color(0xFF1C1B1F),
    surface = Color(0xFFFFFBFE),
    onSurface = Color(0xFF1C1B1F),
    surfaceVariant = Color(0xFFE7E0EC),
    onSurfaceVariant = Color(0xFF49454F),
    outline = Color(0xFF79747E),
    outlineVariant = Color(0xFFCAC4D0),
    surfaceTint = Color(0xFF6750A4),
    inverseSurface = Color(0xFF313033),
    inverseOnSurface = Color(0xFFF4EFF4),
    inversePrimary = Color(0xFFD0BCFF),
)

// ============================================
// DARK COLOR SCHEME (si se necesita)
// ============================================
val DarkColorScheme = darkColorScheme(
    primary = Color(0xFFD0BCFF),
    onPrimary = Color(0xFF381E72),
    primaryContainer = Color(0xFF4F378B),
    onPrimaryContainer = Color(0xFFEADDFF),
    secondary = Color(0xFFCCC2DC),
    onSecondary = Color(0xFF332D41),
    secondaryContainer = Color(0xFF4A4458),
    onSecondaryContainer = Color(0xFFE8DEF8),
    tertiary = Color(0xFFEFB8C8),
    onTertiary = Color(0xFF4A2532),
    tertiaryContainer = Color(0xFF633B48),
    onTertiaryContainer = Color(0xFFFFD8E4),
    error = Color(0xFFFFB4AB),
    onError = Color(0xFF690005),
    errorContainer = Color(0xFF93000A),
    onErrorContainer = Color(0xFFFFDAD6),
    background = Color(0xFF1C1B1F),
    onBackground = Color(0xFFE6E1E5),
    surface = Color(0xFF1C1B1F),
    onSurface = Color(0xFFE6E1E5),
    surfaceVariant = Color(0xFF49454F),
    onSurfaceVariant = Color(0xFFCAC4D0),
    outline = Color(0xFF938F99),
    outlineVariant = Color(0xFF49454F),
    surfaceTint = Color(0xFFD0BCFF),
    inverseSurface = Color(0xFFE6E1E5),
    inverseOnSurface = Color(0xFF313033),
    inversePrimary = Color(0xFF6750A4),
)

// ============================================
// TYPOGRAPHY
// ============================================
val OrienmodTypography = androidx.compose.material3.Typography(
    headlineLarge = HeadlineXL,
    headlineMedium = HeadlineLG,
    headlineSmall = HeadlineLGMobile,
    titleLarge = TitleMD,
    bodyLarge = BodyLG,
    bodyMedium = BodyMD,
    labelLarge = LabelMD,
    labelMedium = LabelSM,
)

// ============================================
// SHAPES
// ============================================
val Shapes = androidx.compose.material3.Shapes(
    // Si no tienes Shapes definidos, usa los por defecto
)

// ============================================
// TEMA PRINCIPAL
// ============================================
@Composable
fun OrienmodTheme(
    darkTheme: Boolean = false,
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) DarkColorScheme else LightColorScheme

    // Configurar la barra de estado para que coincida con el tema
    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as androidx.activity.ComponentActivity).window
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = !darkTheme
        }
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = OrienmodTypography,
        shapes = Shapes,
        content = content
    )
}