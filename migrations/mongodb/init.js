// Inicialización de MongoDB para go_apis_start
// Este script se ejecuta automáticamente cuando Docker inicia MongoDB

// Seleccionar la base de datos
db = db.getSiblingDB('go_apis_db');

// Crear colección de usuarios con índices
db.createCollection('users');
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "is_active": 1 });

// Crear colección de productos con índices
db.createCollection('products');
db.products.createIndex({ "name": 1 });
db.products.createIndex({ "is_active": 1 });
db.products.createIndex({ "price": 1 });

print('✅ Base de datos go_apis_db inicializada correctamente');
print('✅ Colecciones users y products creadas');
print('✅ Índices creados correctamente');
