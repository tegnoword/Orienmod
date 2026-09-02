// com/orienmod/app/theme/Shapes.kt
package com.orienmod.app.theme

import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Shapes
import androidx.compose.ui.unit.dp

val Shapes = Shapes(
    extraSmall = RoundedCornerShape(4.dp),   // sm: 0.25rem
    small = RoundedCornerShape(8.dp),        // DEFAULT: 0.5rem
    medium = RoundedCornerShape(12.dp),      // md: 0.75rem
    large = RoundedCornerShape(16.dp),       // lg: 1rem
    extraLarge = RoundedCornerShape(24.dp),  // xl: 1.5rem
)

// Constantes para uso directo
val ShapeCornerSM = 4.dp
val ShapeCornerDefault = 8.dp
val ShapeCornerMD = 12.dp
val ShapeCornerLG = 16.dp
val ShapeCornerXL = 24.dp
val ShapeCornerFull = 9999.dp

// Extension para usar en componentes
fun Int.toDp() = this.dp