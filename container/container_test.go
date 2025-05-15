package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test interfaces and implementations
type TestService interface {
	GetValue() string
}

type TestServiceImpl struct {
	Value string
}

func (s *TestServiceImpl) GetValue() string {
	return s.Value
}

type DependentService struct {
	Service TestService
}

func (s *DependentService) GetDependencyValue() string {
	return s.Service.GetValue()
}

func TestContainerBasics(t *testing.T) {
	t.Run("Register and resolve singleton", func(t *testing.T) {
		// Arrange
		c := New()
		err := RegisterSingleton[TestService, TestServiceImpl](c)
		require.NoError(t, err)

		// Act
		service1, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Set a value on the service
		service1.(*TestServiceImpl).Value = "test value"

		// Resolve again to verify it's the same instance
		service2, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "test value", service2.GetValue())
		assert.Same(t, service1, service2)
	})

	t.Run("Register and resolve transient", func(t *testing.T) {
		// Arrange
		c := New()
		err := RegisterTransient[TestService, TestServiceImpl](c)
		require.NoError(t, err)

		// Act
		service1, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Set a value on the service
		service1.(*TestServiceImpl).Value = "test value"

		// Resolve again to verify it's a different instance
		service2, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "", service2.GetValue())
		assert.NotSame(t, service1, service2)
	})

	t.Run("Register and resolve scoped", func(t *testing.T) {
		// Arrange
		c := New()
		err := RegisterScoped[TestService, TestServiceImpl](c)
		require.NoError(t, err)

		// Act - Create a scope
		scope := c.CreateScope()

		// Resolve from scope
		service1, err := Resolve[TestService](scope)
		require.NoError(t, err)

		// Set a value on the service
		service1.(*TestServiceImpl).Value = "test value"

		// Resolve again from the same scope to verify it's the same instance
		service2, err := Resolve[TestService](scope)
		require.NoError(t, err)

		// Assert - Same instance within scope
		assert.Equal(t, "test value", service2.GetValue())
		assert.Same(t, service1, service2)

		// Create another scope
		scope2 := c.CreateScope()

		// Resolve from the new scope
		service3, err := Resolve[TestService](scope2)
		require.NoError(t, err)

		// Assert - Different instance in different scope
		assert.Equal(t, "", service3.GetValue())
		assert.NotSame(t, service1, service3)
	})

	t.Run("Register instance", func(t *testing.T) {
		// Arrange
		c := New()
		instance := &TestServiceImpl{Value: "predefined value"}
		err := RegisterInstance[TestService](c, instance)
		require.NoError(t, err)

		// Act
		service, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "predefined value", service.GetValue())
		assert.Same(t, instance, service)
	})

	t.Run("Register factory", func(t *testing.T) {
		// Arrange
		c := New()
		factoryCalled := 0
		RegisterSingletonFactory[TestService](c, func(c *Container) (interface{}, error) {
			factoryCalled++
			return &TestServiceImpl{Value: "factory value"}, nil
		})

		// Act
		service1, err := Resolve[TestService](c)
		require.NoError(t, err)

		service2, err := Resolve[TestService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "factory value", service1.GetValue())
		assert.Equal(t, 1, factoryCalled) // Factory should be called only once for singleton
		assert.Same(t, service1, service2)
	})

	t.Run("Resolve unregistered service", func(t *testing.T) {
		// Arrange
		c := New()

		// Act
		_, err := Resolve[TestService](c)

		// Assert
		assert.Error(t, err)
	})

	t.Run("MustResolve panics for unregistered service", func(t *testing.T) {
		// Arrange
		c := New()

		// Act & Assert
		assert.Panics(t, func() {
			MustResolve[TestService](c)
		})
	})
}

func TestDependencyResolution(t *testing.T) {
	t.Run("Resolve with dependencies", func(t *testing.T) {
		// Arrange
		c := New()

		// Register the dependency
		err := RegisterInstance[TestService](c, &TestServiceImpl{Value: "dependency value"})
		require.NoError(t, err)

		// Register a factory that resolves the dependency
		RegisterSingletonFactory[DependentService](c, func(c *Container) (interface{}, error) {
			dependency, err := Resolve[TestService](c)
			if err != nil {
				return nil, err
			}
			return &DependentService{Service: dependency}, nil
		})

		// Act
		service, err := Resolve[DependentService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "dependency value", service.GetDependencyValue())
	})
}

func TestBuildServiceProvider(t *testing.T) {
	t.Run("Configure multiple services", func(t *testing.T) {
		// Arrange & Act
		c := BuildServiceProvider(func(c *Container) {
			_ = RegisterSingleton[TestService, TestServiceImpl](c)
			RegisterSingletonFactory[DependentService](c, func(c *Container) (interface{}, error) {
				dependency, err := Resolve[TestService](c)
				if err != nil {
					return nil, err
				}
				return &DependentService{Service: dependency}, nil
			})
		})

		// Resolve services
		service1, err := Resolve[TestService](c)
		require.NoError(t, err)
		service1.(*TestServiceImpl).Value = "configured value"

		service2, err := Resolve[DependentService](c)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, "configured value", service2.GetDependencyValue())
	})
}

func TestDispose(t *testing.T) {
	t.Run("Dispose clears services", func(t *testing.T) {
		// Arrange
		c := New()
		err := RegisterSingleton[TestService, TestServiceImpl](c)
		require.NoError(t, err)

		// Verify service can be resolved
		_, err = Resolve[TestService](c)
		require.NoError(t, err)

		// Act
		c.Dispose()

		// Assert
		_, err = Resolve[TestService](c)
		assert.Error(t, err)
	})

	t.Run("Dispose clears scopes", func(t *testing.T) {
		// Arrange
		c := New()
		err := RegisterScoped[TestService, TestServiceImpl](c)
		require.NoError(t, err)

		scope := c.CreateScope()

		// Verify service can be resolved from scope
		_, err = Resolve[TestService](scope)
		require.NoError(t, err)

		// Act
		c.Dispose()

		// Assert - Scope should be disposed
		_, err = Resolve[TestService](scope)
		assert.Error(t, err)
	})
}
