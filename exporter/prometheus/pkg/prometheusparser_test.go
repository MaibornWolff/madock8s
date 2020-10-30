package pkg

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestShouldCreateAListWithMetricBlocks(t *testing.T) {
	Ω := NewGomegaWithT(t)

	metricblocktext := `# HELP jvm_memory_bytes_used Used bytes of a given JVM memory area.
# TYPE jvm_memory_bytes_used gauge
jvm_memory_bytes_used{area="heap",} 1.13428264E8
# HELP jvm_memory_bytes_free Free bytes of a given JVM memory area.
# TYPE jvm_memory_bytes_free counter
jvm_memory_bytes_free{area="heap",} 1.13428264E8`

	metricBlocks := ParsePrometheusResult(metricblocktext)

	expectedMetricItem1 := MetricsMetaBlock{Name: "jvm_memory_bytes_used", Type: "gauge", Help: "Used bytes of a given JVM memory area."}
	expectedMetricItem2 := MetricsMetaBlock{Name: "jvm_memory_bytes_free", Type: "counter", Help: "Free bytes of a given JVM memory area."}

	Ω.Expect(len(metricBlocks)).To(Equal(2))
	Ω.Expect(metricBlocks).Should(ConsistOf(expectedMetricItem1, expectedMetricItem2))
}
