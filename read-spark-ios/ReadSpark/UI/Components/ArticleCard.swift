import SwiftUI

struct ArticleCard: View {
    let article: ArticleSummary
    let onTap: () -> Void

    var body: some View {
        VStack(alignment: .leading, spacing: 8) {
            Text(article.title)
                .font(.headline)

            if let summary = article.summary {
                Text(summary)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
                    .lineLimit(2)
            }

            HStack {
                Label(article.category, systemImage: "folder")
                    .font(.caption)
                Label(article.difficulty, systemImage: "chart.bar")
                    .font(.caption)

                if article.isPremium {
                    Label("VIP", systemImage: "star.fill")
                        .font(.caption)
                        .foregroundColor(.yellow)
                }
            }
        }
        .padding()
        .background(Color(.systemBackground))
        .cornerRadius(12)
        .shadow(radius: 2)
        .onTapGesture(perform: onTap)
    }
}
