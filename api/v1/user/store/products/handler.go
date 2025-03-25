package storeproducts

import (
	"drones-be/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductHandler gestiona la lógica de productos.
type ProductHandler struct {
	db *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// GetAvailableProducts recibe desde el cliente la latitud y longitud
// y retorna los productos que pueden ser entregados, es decir, cuando
// existe una red de warehouses que conectan la ubicación del cliente
// con el proveedor del producto.
// Se asume que:
// - El cliente debe estar a menos de 10 km de algún warehouse.
// - Los drones pueden volar 20 km entre warehouses (para recargarse).
// - Los proveedores están a menos de 10 km de algún warehouse.
func (h *ProductHandler) GetAvailableProducts(c *gin.Context) {
	// Estructura para recibir la posición del cliente.

	_, exist :=c.Get("userID")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	var req struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "detalles": err.Error()})
		return
	}

	// Paso 1: Verificar que el cliente se encuentre en zona de cobertura.
	// Usamos 10km (10000 metros) como radio.
	var count int64
	coverageQuery := `
		SELECT count(*) 
		FROM warehouses 
		WHERE earth_distance(ll_to_earth(?, ?), ll_to_earth(latitude, longitude)) <= 10000
	`
	if err := h.db.Raw(coverageQuery, req.Latitude, req.Longitude).Row().Scan(&count); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar cobertura", "detalles": err.Error()})
		return
	}

	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fuera de la zona de disponibilidad"})
		return
	}

	// Paso 2: Con la siguiente consulta (usando CTE recursivo) se genera la red
	// de warehouses alcanzables a partir de los warehouses que estén a menos de 10 km del cliente.
	// Posteriormente se unen con los proveedores y sus productos.
	//
	// NOTA: Se asume que la tabla "warehouses" tiene los campos "id", "latitude" y "longitude",
	// y que la tabla "providers" tiene "id", "latitude" y "longitude". Además, se usa que un proveedor
	// está conectado a la red si está a menos de 10 km de alguno de los warehouses alcanzables.
	availableQuery := `
	WITH RECURSIVE reachable AS (
		-- Punto de partida: warehouses a menos de 10km del cliente.
		SELECT id, latitude, longitude
		FROM warehouses
		WHERE earth_distance(ll_to_earth(?, ?), ll_to_earth(latitude, longitude)) <= 10000
		
		UNION
		
		-- Se agregan warehouses alcanzables en saltos de hasta 20km.
		SELECT w.id, w.latitude, w.longitude
		FROM warehouses w
		JOIN reachable r ON earth_distance(ll_to_earth(r.latitude, r.longitude), ll_to_earth(w.latitude, w.longitude)) <= 20000
	)
	SELECT DISTINCT prod.id, prod.name, prod.description, prod.price, prod.provider_id
	FROM products prod
	JOIN providers prov ON prov.id = prod.provider_id
	-- Se vincula al warehouse más cercano al proveedor. Se asume que el proveedor
	-- debe estar a menos de 10km de algún warehouse de la red alcanzable.
	JOIN warehouses w ON earth_distance(ll_to_earth(w.latitude, w.longitude), ll_to_earth(prov.latitude, prov.longitude)) <= 10000
	WHERE w.id IN (SELECT id FROM reachable)
	`

	var products []models.Product
	if err := h.db.Raw(availableQuery, req.Latitude, req.Longitude).Scan(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos", "detalles": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
